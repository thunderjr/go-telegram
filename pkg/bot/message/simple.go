package message

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ message = (*SimpleMessage)(nil)

type SimpleMessage struct {
	sender    sender
	Content   string
	Recipient int64
	MessageID int

	OnReply string
}

type Params struct {
	Bot       sender
	Content   string
	Recipient int64
	OnReply   string
}

func NewSimpleMessage(p *Params) *SimpleMessage {
	return &SimpleMessage{
		Recipient: p.Recipient,
		Content:   p.Content,
		OnReply:   p.OnReply,
		sender:    p.Bot,
	}
}

func (s *SimpleMessage) Send(ctx context.Context, opts ...sendOption) (*tgbotapi.Message, error) {
	var msg tgbotapi.Chattable = tgbotapi.NewMessage(s.Recipient, s.Content)

	for _, opt := range opts {
		msg = opt(msg)
	}

	m, err := s.sender.Send(msg)
	if err != nil {
		return nil, err
	}

	s.MessageID = m.MessageID
	if s.OnReply != "" {
		if err := ReplyActionRepo(ctx).Save(ctx, newReplyActionFromSimple(s)); err != nil {
			log.Printf("Error saving message reply action: %v\n", err)
		}
	}

	return &m, nil
}

func (s SimpleMessage) GetID() string {
	if s.OnReply != "" {
		return fmt.Sprintf("%d:%d", s.Recipient, s.MessageID)
	}
	return ""
}
