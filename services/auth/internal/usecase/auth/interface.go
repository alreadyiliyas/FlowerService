package usecase

import (
	"context"

	"github.com/ilyas/flower/services/auth/internal/dto"
)

type UsecaseAuth interface {
	Login(ctx context.Context, dtoReq dto.LoginRequest) (dtoRes *dto.LoginResponse, err error)
	Registration(ctx context.Context, dtoReq dto.RegistrationRequest) (dtoRes *dto.RegistrationResponse, err error)
	VerifyAccount(ctx context.Context, dtoReq dto.VerifyAccountRequest) error
	SetPassword(ctx context.Context, dtoReq dto.SetPasswordRequest) error
	SendCodeToUpdatePassword(ctx context.Context, phoneNumber *string) error
	UpdatePassword(ctx context.Context, dtoReq dto.ConfirmUpdatePasswordRequest) error
	RefreshToken(ctx context.Context, dtoReq dto.RefreshTokenRequest) (dtoRes *dto.RefreshTokenResponse, err error)
	Logout(ctx context.Context, dtoReq dto.LogoutRequest) error
	LogoutAll(ctx context.Context, dtoReq dto.LogoutAllRequest) error
	IsSessionActive(ctx context.Context, sessionID string) (bool, error)
}
