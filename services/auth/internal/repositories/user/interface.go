package repositories

import (
	"context"

	"github.com/ilyas/flower/services/auth/internal/entities"
)

type UserRepository interface {
	Get(ctx context.Context, account *entities.Auth) (*entities.User, error)
	Update(ctx context.Context, account *entities.Auth) (*entities.User, error)
	Delete(ctx context.Context, account *entities.Auth) error
}
