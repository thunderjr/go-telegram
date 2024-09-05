package message

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendOption func(tgbotapi.Chattable) tgbotapi.Chattable

type message interface {
	Send(ctx context.Context, opts ...sendOption) (*tgbotapi.Message, error)
}

type sender interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}
