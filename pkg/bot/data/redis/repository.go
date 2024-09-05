package redis

import (
	"context"
	"encoding/json"

	"github.com/thunderjr/go-telegram/pkg/bot/data"
)

func (r *RedisRepository[T]) FindOne(ctx context.Context, query T) (*T, error) {
	raw, err := r.client.Get(ctx, r.getPrefixed(query.GetID())).Result()
	if err != nil {
		return nil, data.ErrInternal(err)
	}

	if len(raw) == 0 {
		return nil, data.ErrNotFound
	}

	res := new(T)
	if err := json.Unmarshal([]byte(raw), res); err != nil {
		return nil, data.ErrInternal(err)
	}

	return res, nil
}

func (r *RedisRepository[T]) Save(ctx context.Context, d T) error {
	raw, err := json.Marshal(d)
	if err != nil {
		return data.ErrInternal(err)
	}

	if err := r.client.Set(ctx, r.getPrefixed(d.GetID()), raw, r.ttl).Err(); err != nil {
		return data.ErrInternal(err)
	}

	return nil
}

func (r *RedisRepository[T]) Remove(ctx context.Context, query T) error {
	if _, err := r.client.Del(ctx, r.getPrefixed(query.GetID())).Result(); err != nil {
		return data.ErrInternal(err)
	}
	return nil
}
