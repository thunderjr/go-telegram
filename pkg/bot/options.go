package bot

import "github.com/thunderjr/go-telegram/pkg/bot/update"

type botOption func(*TelegramBot)

func WithUpdateHandlers(handlers []update.Handler) botOption {
	return func(b *TelegramBot) {
		b.UpdateGateway = update.NewGateway(handlers...)
	}
}
