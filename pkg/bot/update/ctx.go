package update

import (
	"context"

	"github.com/thunderjr/go-telegram/pkg/bot/data"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
)

func ReplyActionRepo(ctx context.Context) data.Repository[message.ReplyAction] {
	r, ok := ctx.Value("ReplyActionRepo").(data.Repository[message.ReplyAction])
	if !ok {
		panic("context key ReplyActionRepo not found")
	}
	return r
}

func FormAnswerRepo(ctx context.Context) data.Repository[FormAnswer] {
	r, ok := ctx.Value("FormAnswerRepo").(data.Repository[FormAnswer])
	if !ok {
		panic("context key FormAnswerRepo not found")
	}
	return r
}
