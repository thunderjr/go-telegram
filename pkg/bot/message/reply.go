package message

import "fmt"

type ReplyAction struct {
	MessageID int
	Recipient int64
	OnReply   string
}

func (r ReplyAction) GetID() string {
	return fmt.Sprintf("%d:%d", r.Recipient, r.MessageID)
}

func newReplyActionFromSimple(m *SimpleMessage) ReplyAction {
	return ReplyAction{
		MessageID: m.MessageID,
		Recipient: m.Recipient,
		OnReply:   m.OnReply,
	}
}

func newReplyActionFromCandidate(c *CandidateMessage) ReplyAction {
	return ReplyAction{
		MessageID: c.MessageID,
		Recipient: c.Recipient,
		OnReply:   c.OnReply,
	}
}
