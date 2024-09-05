package message

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const editableChatIDPattern = "editable_msg:%d"

var _ message = (*EditableMessage)(nil)

type EditableMessage struct {
	msg       message
	Key       string
	MessageID int
}

func ToEditable(msg message) *EditableMessage {
	return &EditableMessage{msg: msg}
}

func (e *EditableMessage) Send(ctx context.Context, opts ...sendOption) (*tgbotapi.Message, error) {
	msg, err := e.msg.Send(ctx, opts...)
	if err != nil {
		return nil, err
	}

	e.MessageID = msg.MessageID
	e.Key = getEditableID(msg.Chat.ID)
	if err := EditableRepo(ctx).Save(ctx, *e); err != nil {
		return nil, err
	}

	return msg, nil
}

func (e EditableMessage) GetID() string {
	return e.Key
}

func getEditableID(chatID int64) string {
	return fmt.Sprintf(editableChatIDPattern, chatID)
}
