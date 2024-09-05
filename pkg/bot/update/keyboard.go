package update

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var _ Handler = (*KeyboardCallbackUpdate)(nil)

type KeyboardCallbackUpdate struct {
	fn  func(tgbotapi.Update) error
	key string
}

func NewKeyboardCallbackUpdate(key string, fn func(tgbotapi.Update) error) *KeyboardCallbackUpdate {
	return &KeyboardCallbackUpdate{fn, key}
}

func (KeyboardCallbackUpdate) Type() HandlerType {
	return HandlerTypeKeyboardCallback
}

func (k KeyboardCallbackUpdate) Key() string {
	return k.key
}

func (k KeyboardCallbackUpdate) Handle(u tgbotapi.Update) error {
	return k.fn(u)
}
