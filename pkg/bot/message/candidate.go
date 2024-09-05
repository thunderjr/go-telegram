package message

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ message = (*CandidateMessage)(nil)

type CandidateMessage struct {
	sender    sender
	Content   string
	OnReply   string
	Recipient int64
	MessageID int
}

func NewCandidateMessage(p *Params) *CandidateMessage {
	return &CandidateMessage{
		Recipient: p.Recipient,
		Content:   p.Content,
		OnReply:   p.OnReply,
		sender:    p.Bot,
	}
}

func (c *CandidateMessage) Send(ctx context.Context, opts ...sendOption) (*tgbotapi.Message, error) {
	var msg tgbotapi.Chattable = tgbotapi.NewMessage(c.Recipient, c.Content)

	lastEditable, err := EditableRepo(ctx).FindOne(ctx, EditableMessage{
		Key: getEditableID(c.Recipient),
	})
	if err == nil {
		msg = tgbotapi.NewEditMessageText(c.Recipient, lastEditable.MessageID, c.Content)
		defer func() {
			if err := EditableRepo(ctx).Remove(ctx, *lastEditable); err != nil {
				log.Printf("failed to remove editable message: %v\n", err)
			}
		}()
	}

	for _, opt := range opts {
		msg = opt(msg)
	}

	res, err := c.sender.Send(msg)
	if err != nil {
		return nil, err
	}

	c.MessageID = res.MessageID
	if c.OnReply != "" {
		if err := ReplyActionRepo(ctx).Save(ctx, newReplyActionFromCandidate(c)); err != nil {
			log.Printf("Error saving message reply action: %v\n", err)
		}
	}

	return &res, nil
}
