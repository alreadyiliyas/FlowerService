package repositories

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) CacheRepository {
	return &redisRepo{client: client}
}

func (r *redisRepo) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisRepo) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisRepo) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
