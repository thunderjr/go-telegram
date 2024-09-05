package update

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var _ Handler = (*WebappDataUpdate)(nil)

type WebappDataUpdate struct {
	fn  func(tgbotapi.Update) error
	key string
}

func NewWebappDataUpdate(key string, fn func(tgbotapi.Update) error) *WebappDataUpdate {
	return &WebappDataUpdate{fn, key}
}

func (WebappDataUpdate) Type() HandlerType {
	return HandlerTypeWebApp
}

func (w WebappDataUpdate) Key() string {
	return w.key
}

func (w WebappDataUpdate) Handle(u tgbotapi.Update) error {
	return w.fn(u)
}
