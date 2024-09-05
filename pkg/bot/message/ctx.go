package message

import (
	"context"

	"github.com/thunderjr/go-telegram/pkg/bot/data"
)

func EditableRepo(ctx context.Context) data.Repository[EditableMessage] {
	r, ok := ctx.Value("EditableRepo").(data.Repository[EditableMessage])
	if !ok {
		panic("context key EditableRepo not found")
	}
	return r
}

func ReplyActionRepo(ctx context.Context) data.Repository[ReplyAction] {
	r, ok := ctx.Value("ReplyActionRepo").(data.Repository[ReplyAction])
	if !ok {
		panic("context key ReplyActionRepo not found")
	}
	return r
}
