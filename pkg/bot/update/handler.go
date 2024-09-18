package update

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var _ Handler = (*updateHandler)(nil)

type updateHandler struct {
	fn  func(tgbotapi.Update) error
	t   HandlerType
	key string
}

func newUpdateHandler(
	t HandlerType,
	key string,
	fn func(tgbotapi.Update) error,
) *updateHandler {
	return &updateHandler{fn, t, key}
}

func (m updateHandler) Type() HandlerType {
	return m.t
}

func (m updateHandler) Key() string {
	return m.key
}

func (m updateHandler) Handle(u tgbotapi.Update) error {
	return m.fn(u)
}
