package update

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
)

var _ Handler = (*ReplyUpdate)(nil)

type ReplyUpdate struct {
	fn  func(tgbotapi.Update) error
	key string
}

func NewReplyUpdate(key string, fn func(tgbotapi.Update) error) *ReplyUpdate {
	return &ReplyUpdate{fn, key}
}

func (ReplyUpdate) Type() HandlerType {
	return HandlerTypeReply
}

func (m ReplyUpdate) Key() string {
	return m.key
}

func (m ReplyUpdate) Handle(u tgbotapi.Update) error {
	return m.fn(u)
}

func getReplyAction(ctx context.Context, u tgbotapi.Update) (string, error) {
	m, err := ReplyActionRepo(ctx).FindOne(ctx, message.ReplyAction{
		MessageID: u.Message.ReplyToMessage.MessageID,
		Recipient: u.Message.Chat.ID,
	})
	if err != nil {
		return "", err
	}
	return m.OnReply, nil
}
