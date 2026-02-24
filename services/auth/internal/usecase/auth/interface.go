package usecase

import (
	"context"

	"github.com/ilyas/flower/services/auth/internal/dto"
)

type Usecase interface {
	// Login(ctx context.Context, dtoReq dto.LoginRequest) (dtoRes *dto.LoginResponse, err error)
	Registration(ctx context.Context, dtoReq dto.RegistrationRequest) (dtoRes *dto.RegistrationResponse, err error)
	// VerifyAccount(ctx context.Context, dtoReq dto.VerifyAccountRequest) error
	// SetPassword(ctx context.Context, dtoReq dto.SetPasswordRequest) error
	// UpdatePassword(ctx context.Context, dtoReq dto.UpdatePasswordRequest) error
	// RefreshToken(ctx context.Context, dtoReq dto.RefreshTokenRequest) (dtoRes *dto.RefreshTokenResponse, err error)
	// Logout(ctx context.Context, dtoReq dto.LogoutRequest) error
}
