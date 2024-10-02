package redisdb

import (
	"github.com/bccfilkom/career-path-service/internal/pkg/env"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

func NewInstance() (*redis.Client, error) {
	DBHost := env.GetString("REDIS_HOST", "")
	DBPort := env.GetString("REDIS_PORT", "")
	DBPassword := env.GetString("REDIS_PASSWORD", "")

	options := &redis.Options{
		Addr:     DBHost + ":" + DBPort,
		Password: DBPassword,
		DB:       0,
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
