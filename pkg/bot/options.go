package bot

import "github.com/thunderjr/go-telegram/pkg/bot/update"

type BotOption func(*TelegramBot)

func WithUpdateHandlers(handlers []update.Handler) BotOption {
	return func(b *TelegramBot) {
		b.updateGateway = update.NewGateway(handlers...)
	}
}
