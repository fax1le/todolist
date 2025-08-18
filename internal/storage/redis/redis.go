package redis

import (
	"context"
	"todo/internal/log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Client     redis.Client
	REDIS_HOST = os.Getenv("REDIS_HOST")
)


func StartRedis() {
	Client = *redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	err := Client.Ping(context.Background()).Err()

	if err != nil {
		log.Logger.Error("Redis connection failed", "err", err)
		os.Exit(1)
	}

	log.Logger.Info("Redis connection established")
}
