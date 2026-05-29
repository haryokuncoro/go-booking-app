package cache

import (
	"context"
	"fmt"

	"booking-app/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(
	cfg *config.Config,
) *redis.Client {

	client := redis.NewClient(
		&redis.Options{
			Addr: fmt.Sprintf(
				"%s:%s",
				cfg.RedisHost,
				cfg.RedisPort,
			),
		},
	)

	client.Ping(
		context.Background(),
	)

	return client
}