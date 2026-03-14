package usecase

import (
	"context"

	"github.com/ilyas/flower/services/auth/internal/dto"
)

type UsecaseUser interface {
	GetUserInfo(ctx context.Context, dtoReq dto.GetUserInfoRequest) (dtoRes *dto.GetUserInfoResponse, err error)
	UpdateUserInfo(ctx context.Context, dtoReq dto.UpdateUserInfoRequest) (dtoRes *dto.UpdateUserInfoResponse, err error)
	DeleteUserInfo(ctx context.Context, dtoReq dto.DeleteUserRequest) error
}
