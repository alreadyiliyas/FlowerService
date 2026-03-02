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

func (r *redisRepo) SaveConfirmationCode(ctx context.Context, phone, code *string, ttl time.Duration) error {
	return r.client.Set(ctx, *phone, *code, ttl).Err()
}

func (r *redisRepo) GetConfirmCode(ctx context.Context, phone *string) (string, error) {
	var code string
	err := r.client.Get(ctx, *phone).Scan(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (r *redisRepo) DeleteConfirmCode(ctx context.Context, phone *string) error {
	return r.client.Del(ctx, *phone).Err()
}
