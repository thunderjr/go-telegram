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
