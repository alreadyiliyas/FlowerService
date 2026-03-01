package repositories

import (
	"context"
	"time"
)

type CacheRepository interface {
	SaveConfirmationCode(ctx context.Context, phone, code *string, ttl time.Duration) error
	GetConfirmCode(ctx context.Context, phone *string) (string, error)
}
