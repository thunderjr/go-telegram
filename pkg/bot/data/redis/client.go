package redis

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

var instance *redis.Client

func Instance() *redis.Client {
	var once sync.Once
	once.Do(func() {
		password := os.Getenv("REDIS_PASSWORD")
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			host = "127.0.0.1"
		}

		port := 6379
		if p, err := strconv.Atoi(os.Getenv("REDIS_PORT")); err == nil {
			port = p
		}

		db := 0
		if d, err := strconv.Atoi(os.Getenv("REDIS_DB")); err == nil {
			db = d
		}

		instance = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: password,
			DB:       db,
		})
	})
	return instance
}
