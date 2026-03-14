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

func (r *redisRepo) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *redisRepo) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

func (r *redisRepo) SAdd(ctx context.Context, key string, members ...string) error {
	args := make([]interface{}, 0, len(members))
	for _, m := range members {
		args = append(args, m)
	}
	return r.client.SAdd(ctx, key, args...).Err()
}

func (r *redisRepo) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

func (r *redisRepo) SRem(ctx context.Context, key string, members ...string) error {
	args := make([]interface{}, 0, len(members))
	for _, m := range members {
		args = append(args, m)
	}
	return r.client.SRem(ctx, key, args...).Err()
}
