package update

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var _ Handler = (*MessageUpdate)(nil)

type MessageUpdate struct {
	fn  func(tgbotapi.Update) error
	key string
}

func NewMessageUpdate(key string, fn func(tgbotapi.Update) error) *MessageUpdate {
	return &MessageUpdate{fn, key}
}

func (MessageUpdate) Type() HandlerType {
	return HandlerTypeMessage
}

func (m MessageUpdate) Key() string {
	return m.key
}

func (m MessageUpdate) Handle(u tgbotapi.Update) error {
	return m.fn(u)
}
