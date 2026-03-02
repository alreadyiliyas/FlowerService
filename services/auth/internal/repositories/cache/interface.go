package repositories

import (
	"context"
	"time"
)

type CacheRepository interface {
	SaveConfirmationCode(ctx context.Context, phone, code *string, ttl time.Duration) error
	DeleteConfirmCode(ctx context.Context, phone *string) error
	GetConfirmCode(ctx context.Context, phone *string) (string, error)
}
