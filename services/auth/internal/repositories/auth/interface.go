package repositories

import (
	"context"

	"github.com/ilyas/flower/services/auth/internal/entities"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *entities.User, account *entities.Auth) (*entities.User, error)
	GetAccountByPhoneNumber(ctx context.Context, in *string) (*entities.Auth, error)
	VerifyAccount(ctx context.Context, account *entities.Auth) error
	SetPassword(ctx context.Context, account *entities.Auth) error
	GetPassword(ctx context.Context, in *string) (*string, error)
}
