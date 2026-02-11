package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/savanyv/bsnack-backend/config"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB: 0,
	})
	return &RedisClient{
		Client: rdb,
	}
}

func (r *RedisClient) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
