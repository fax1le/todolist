package redis

import (
	"context"
	"log/slog"
	"todo/internal/config"

	"github.com/redis/go-redis/v9"
)

func StartRedis(cfg config.Config, logger *slog.Logger) (*redis.Client, error) {
	Client := redis.NewClient(&redis.Options{
		Addr:     "redis:" + cfg.RedisHost,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDb,
		Protocol: cfg.RedisProtocol,
	})

	err := Client.Ping(context.Background()).Err()

	return Client, err
}
