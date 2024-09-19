package update

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thunderjr/go-telegram/pkg/bot/data"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
)

type formStep struct {
	FieldName string
	Prompt    string
	Action    string // "prompt" or "submit"
}

type FormAnswer struct {
	Steps        []formStep
	Key          string
	UserID       int64
	CurrentIndex int
	Form         any
}

type FormFieldPrompt struct {
	Name   string
	Prompt string
	Order  int
}

type NewFormHandlerParams[T any] struct {
	OnSubmit func(ctx context.Context, data *T) error
	Bot      interface {
		Send(tgbotapi.Chattable) (tgbotapi.Message, error)
	}

	Type HandlerType
	Key  string
	Form *T
}

func (f FormAnswer) GetID() string {
	return fmt.Sprintf("form:%d:%s", f.UserID, f.Key)
}

func NewFormHandlers[T any](
	ctx context.Context,
	params *NewFormHandlerParams[T],
) []Handler {
	repo := FormAnswerRepo(ctx)
	replyKey := fmt.Sprintf("reply:%s", params.Key)

	startHandler := newUpdateHandler(params.Type, params.Key, func(update tgbotapi.Update) error {
		userID := getChatID(update)

		prompts, err := promptFromFields(params.Form)
		if err != nil {
			return err
		}

		if len(prompts) == 0 {
			return fmt.Errorf("no prompts found for form %s", params.Key)
		}

		slices.SortFunc(prompts, func(i, j FormFieldPrompt) int {
			return i.Order - j.Order
		})

		formSteps := make([]formStep, len(prompts))
		for i, prompt := range prompts {
			action := "prompt"
			if i == len(prompts)-1 {
				action = "submit"
			}

			formSteps[i] = formStep{
				Prompt:    prompt.Prompt,
				FieldName: prompt.Name,
				Action:    action,
			}
		}

		answer := &FormAnswer{
			UserID: userID,
			Steps:  formSteps,
			Key:    params.Key,
			Form:   params.Form,
		}

		firstStep := formSteps[answer.CurrentIndex]
		msg := message.NewSimpleMessage(&message.Params{
			Content:   firstStep.Prompt,
			Bot:       params.Bot,
			OnReply:   replyKey,
			Recipient: userID,
		})
		if _, err := msg.Send(ctx, message.WithForceReply()); err != nil {
			return err
		}

		return repo.Save(ctx, *answer)
	})

	updateHandler := newUpdateHandler(HandlerTypeReply, replyKey, func(update tgbotapi.Update) error {
		userID := update.Message.Chat.ID

		answer, err := repo.FindOne(ctx, FormAnswer{UserID: userID, Key: params.Key})
		if err != nil && err != data.ErrNotFound {
			return err
		}

		if answer == nil {
			return fmt.Errorf("no form in progress for user %d", userID)
		}

		if answer.Form == nil {
			return fmt.Errorf("form answer has no form data")
		}

		jsonForm, err := json.Marshal(answer.Form)
		if err != nil {
			return fmt.Errorf("error marshalling form data: %v", err)
		}

		form := new(T)
		if err := json.Unmarshal(jsonForm, &form); err != nil {
			return fmt.Errorf("error unmarshalling form data: %v", err)
		}

		step := answer.Steps[answer.CurrentIndex]
		fieldName := step.FieldName

		if err := setField(form, fieldName, update.Message.Text); err != nil {
			return err
		}

		answer.Form = form

		if step.Action == "submit" {
			if params.OnSubmit != nil {
				if err := params.OnSubmit(ctx, form); err != nil {
					log.Printf("error submitting form answer (%v): %v\n", answer, err)
					return err
				}
			}

			if err := repo.Remove(ctx, *answer); err != nil {
				log.Printf("error removing form answer (%v): %v\n", answer, err)
				return err
			}

			return nil
		}

		if answer.CurrentIndex+1 >= len(answer.Steps) {
			return nil
		}

		answer.CurrentIndex++
		step = answer.Steps[answer.CurrentIndex]

		msg := message.NewSimpleMessage(&message.Params{
			Content:   step.Prompt,
			Bot:       params.Bot,
			OnReply:   replyKey,
			Recipient: userID,
		})
		if _, err := msg.Send(ctx, message.WithForceReply()); err != nil {
			return err
		}

		if err := repo.Save(ctx, *answer); err != nil {
			return err
		}

		return nil
	})

	return []Handler{startHandler, updateHandler}
}

// PromptProvider is an interface that should be implemented by structs
// that provide custom prompts for their fields.
type PromptProvider interface {
	FieldPrompts() ([]FormFieldPrompt, error)
}

func promptFromFields(s interface{}) (prompts []FormFieldPrompt, err error) {
	if provider, ok := s.(PromptProvider); ok {
		prompts, err = provider.FieldPrompts()
		if err != nil {
			return nil, err
		}
	}

	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("promptFromFields: expected a pointer to a struct")
	}

	el := val.Elem().Type()
	lf := el.NumField()
	for i := 0; i < lf; i++ {
		field := el.Field(i)
		prompt := field.Tag.Get("telegram_prompt")
		if prompt != "" {
			order, err := strconv.Atoi(field.Tag.Get("telegram_prompt_order"))
			if err != nil {
				order = i
			}

			prompts = append(prompts, FormFieldPrompt{
				Name:   field.Name,
				Prompt: prompt,
				Order:  order,
			})
		}
	}

	return prompts, nil
}

func setField(s interface{}, fieldName, value string) error {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("setField: expected a pointer to a struct")
	}

	field := val.Elem().FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("setField: no such field: %s in struct", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("setField: cannot set field %s", fieldName)
	}

	switch field.Type().Kind() {
	case reflect.Int:
		iv, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("error converting value to int: (%s) %w", value, err)
		}
		field.SetInt(int64(iv))
	case reflect.Float64:
		fv, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("error converting value to float64: (%s) %w", value, err)
		}
		field.SetFloat(fv)
	case reflect.String:
		field.SetString(value)
	default:
		fmt.Println("[WARNING] This type cast is not yet implemented.")
	}
	return nil
}

func getChatID(update tgbotapi.Update) int64 {
	if update.Message == nil {
		return update.CallbackQuery.Message.Chat.ID
	}
	return update.Message.Chat.ID
}
