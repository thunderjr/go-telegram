package redis

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thunderjr/go-telegram/pkg/bot/data"
)

type RedisRepository[T data.Entity] struct {
	client *redis.Client
	ttl    time.Duration
	prefix string
}

type Config struct {
	Client *redis.Client
	// TODO: maybe create an option interface to custom repository configuration
	TTL    time.Duration
	Prefix string
}

func NewRepository[T data.Entity](cfg *Config) data.Repository[T] {
	if cfg.Client != nil {
		return &RedisRepository[T]{cfg.Client, cfg.TTL, cfg.Prefix}
	}
	client := Instance()
	return &RedisRepository[T]{client, cfg.TTL, cfg.Prefix}
}

func (r RedisRepository[T]) getPrefixed(id string) string {
	return fmt.Sprintf("%s:%s", r.prefix, id)
}
