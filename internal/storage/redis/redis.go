package redis

import (
	"context"
	"log/slog"
	"os"
	"todo/internal/config"

	"github.com/redis/go-redis/v9"
)

func StartRedis(cfg config.Config, logger *slog.Logger) *redis.Client {
	Client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDb,
		Protocol: cfg.RedisProtocol,
	})

	err := Client.Ping(context.Background()).Err()

	if err != nil {
		logger.Error("redis connection failed", "err", err)
		os.Exit(1)
	}

	logger.Info("redis connection established")
	return Client
}
