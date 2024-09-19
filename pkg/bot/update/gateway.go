package update

import (
	"context"
	"errors"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
)

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

func (g *Gateway) AddHandlers(handlers ...Handler) {
	for _, handler := range handlers {
		g.handlers[handler.Key()] = handler
		g.count++
	}
}

func (g Gateway) Handle(ctx context.Context, update tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		for key, handler := range g.handlers {
			if strings.HasPrefix(update.CallbackQuery.Data, key) && handler.Type() == HandlerTypeKeyboardCallback {
				return handler.Handle(update)
			}
		}
	}

	if update.Message != nil {
		if update.Message.ReplyToMessage != nil {
			action, err, remove := getReplyAction(ctx, update)
			if err != nil {
				return errors.Join(err, remove())
			}

			handler, ok := g.handlers[action]
			if ok && handler.Type() == HandlerTypeReply {
				return firstError(handler.Handle(update), remove())
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

func getReplyAction(ctx context.Context, u tgbotapi.Update) (string, error, func() error) {
	removeFn := func() error {
		return ReplyActionRepo(ctx).Remove(ctx, message.ReplyAction{
			MessageID: u.Message.ReplyToMessage.MessageID,
			Recipient: u.Message.Chat.ID,
		})
	}

	m, err := ReplyActionRepo(ctx).FindOne(ctx, message.ReplyAction{
		MessageID: u.Message.ReplyToMessage.MessageID,
		Recipient: u.Message.Chat.ID,
	})
	if err != nil {
		return "", err, removeFn
	}
	return m.OnReply, nil, removeFn
}

func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
