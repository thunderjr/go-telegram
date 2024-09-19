package bot

import (
	"context"

	"github.com/thunderjr/go-telegram/pkg/bot/data"
	"github.com/thunderjr/go-telegram/pkg/bot/data/redis"
	"github.com/thunderjr/go-telegram/pkg/bot/message"
	"github.com/thunderjr/go-telegram/pkg/bot/update"
)

type dataProvider string

var (
	RedisProvider dataProvider = "redis"
)

var (
	AppName      = ""
	DataProvider = RedisProvider
)

func NewRepository[T data.Entity]() data.Repository[T] {
	if AppName == "" {
		panic("App name is required")
	}

	switch DataProvider {
	case RedisProvider:
		return redis.NewRepository[T](&redis.Config{
			Prefix: "telegram:" + AppName,
		})
	default:
		return redis.NewRepository[T](&redis.Config{
			Prefix: "telegram:" + AppName,
		})
	}
}

func SetAppName(name string) {
	AppName = name
}

func SetDataProvider(provider dataProvider) {
	DataProvider = provider
}

func WithEditableRepo(ctx context.Context, repo data.Repository[message.EditableMessage]) context.Context {
	return context.WithValue(ctx, "EditableRepo", repo)
}

func WithReplyActionRepo(ctx context.Context, repo data.Repository[message.ReplyAction]) context.Context {
	return context.WithValue(ctx, "ReplyActionRepo", repo)
}

func WithFormAnswerRepo(ctx context.Context, repo data.Repository[update.FormAnswer]) context.Context {
	return context.WithValue(ctx, "FormAnswerRepo", repo)
}
