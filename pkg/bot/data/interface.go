package data

import (
	"context"
)

type (
	Repository[T Entity] interface {
		FindOne(ctx context.Context, query T) (*T, error)
		Remove(ctx context.Context, query T) error
		Save(ctx context.Context, data T) error
	}

	Entity interface {
		GetID() string
	}
)
