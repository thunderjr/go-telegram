package update

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	HandlerTypeKeyboardCallback HandlerType = "keyboard_callback"
	HandlerTypeMessage          HandlerType = "message"
	HandlerTypeWebApp           HandlerType = "webapp"
	HandlerTypeReply            HandlerType = "reply"
)

type HandlerType string

type Handler interface {
	Key() string
	Type() HandlerType
	Handle(tgbotapi.Update) error
}

type Gateway struct {
	handlers map[string]Handler
	count    int
}

func NewGateway(handlers ...Handler) *Gateway {
	var count int
	h := make(map[string]Handler)
	for _, handler := range handlers {
		h[handler.Key()] = handler
		count++
	}
	return &Gateway{h, count}
}

func (g Gateway) Handle(ctx context.Context, update tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		handler, ok := g.handlers[update.CallbackQuery.Data]
		if ok && handler.Type() == HandlerTypeKeyboardCallback {
			return handler.Handle(update)
		}
	}

	if update.Message != nil {
		if update.Message.ReplyToMessage != nil {
			action, err := getReplyAction(ctx, update)
			if err != nil {
				return err
			}

			fmt.Println("action found", action)
			handler, ok := g.handlers[action]
			fmt.Println("action ok", ok)
			if ok && handler.Type() == HandlerTypeReply {
				fmt.Println("Calling this shit")
				return handler.Handle(update)
			}
		}

		if update.Message.WebAppData != nil {
			handler, ok := g.handlers[update.Message.WebAppData.ButtonText]
			if ok && handler.Type() == HandlerTypeWebApp {
				return handler.Handle(update)
			}
		}

		handler, ok := g.handlers[strings.TrimPrefix(update.Message.Text, "/")]
		if ok && handler.Type() == HandlerTypeMessage {
			return handler.Handle(update)
		}
	}

	return nil
}

func (g Gateway) Len() int {
	return g.count
}
