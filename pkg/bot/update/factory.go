package update

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	HandlerTypeKeyboardCallback HandlerType = "keyboard_callback"
	HandlerTypeMessage          HandlerType = "message"
	HandlerTypeWebApp           HandlerType = "webapp"
	HandlerTypeReply            HandlerType = "reply"
)

type HandlerType string

type Handler interface {
	Key() string
	Type() HandlerType
	Handle(tgbotapi.Update) error
}

func NewKeyboardCallbackUpdate(key string, fn func(tgbotapi.Update) error) Handler {
	return newUpdateHandler(HandlerTypeKeyboardCallback, key, fn)
}

func NewMessageUpdate(key string, fn func(tgbotapi.Update) error) Handler {
	return newUpdateHandler(HandlerTypeMessage, key, fn)
}

func NewWebappDataUpdate(key string, fn func(tgbotapi.Update) error) Handler {
	return newUpdateHandler(HandlerTypeWebApp, key, fn)
}

func NewReplyUpdate(key string, fn func(tgbotapi.Update) error) Handler {
	return newUpdateHandler(HandlerTypeReply, key, fn)
}
