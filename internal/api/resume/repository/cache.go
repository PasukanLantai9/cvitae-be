package resumeRepository

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil // Key does not exist
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

// Set stores a value in Redis with a specified expiration time
func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a key from Redis
func (r *redisRepository) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
