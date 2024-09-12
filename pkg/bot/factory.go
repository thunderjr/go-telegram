package bot

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thunderjr/go-telegram/pkg/bot/update"
)

type TelegramBot struct {
	*tgbotapi.BotAPI
	UpdateGateway *update.Gateway
}

func New(token string, opts ...BotOption) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	instance := &TelegramBot{BotAPI: bot}
	for _, opt := range opts {
		opt(instance)
	}

	return instance, nil
}

func (t *TelegramBot) Updates(ctx context.Context, errChan chan<- error) {
	if t.UpdateGateway == nil {
		log.Println("[update gateway] not initialized")
		return
	}

	if t.UpdateGateway.Len() == 0 {
		log.Println("[update gateway] no handlers")
		return
	}

	for update := range t.GetUpdatesChan(tgbotapi.NewUpdate(0)) {
		if err := t.UpdateGateway.Handle(ctx, update); err != nil {
			errChan <- fmt.Errorf("[update gateway] error handling update: %w", err)
			continue
		}
	}
}
