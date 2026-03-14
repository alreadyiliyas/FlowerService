package usecase

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ilyas/flower/services/auth/internal/apperrors"
	"github.com/ilyas/flower/services/auth/internal/config"
	"github.com/ilyas/flower/services/auth/internal/dto"
	"github.com/ilyas/flower/services/auth/internal/entities"
	cache "github.com/ilyas/flower/services/auth/internal/repositories/cache"
	user "github.com/ilyas/flower/services/auth/internal/repositories/user"
	"github.com/ilyas/flower/services/auth/internal/utils"
)

type userUsecase struct {
	cfg       *config.Config
	trRepo    user.UserRepository
	cacheRepo cache.CacheRepository
}

func New(cfg *config.Config, ur user.UserRepository, cr cache.CacheRepository) UsecaseUser {
	return &userUsecase{
		cfg:       cfg,
		trRepo:    ur,
		cacheRepo: cr,
	}
}

func (uc *userUsecase) GetUserInfo(ctx context.Context, dtoReq dto.GetUserInfoRequest) (dtoRes *dto.GetUserInfoResponse, err error) {
	if dtoReq.PhoneNumber == nil || *dtoReq.PhoneNumber == "" {
		return nil, fmt.Errorf("%w: задан пустой номер телефона", apperrors.ErrInvalidInput)
	}
	if !utils.IsValidPhoneNumber(*dtoReq.PhoneNumber) {
		return nil, fmt.Errorf("%w: проверьте номер телефона", apperrors.ErrInvalidInput)
	}

	cacheKey := utils.BuildUserInfoKey(*dtoReq.PhoneNumber)
	if raw, err := uc.cacheRepo.Get(ctx, cacheKey); err == nil {
		var cached dto.GetUserInfoResponse
		if err := utils.UnmarshalFromString(raw, &cached); err != nil {
			log.Printf("| usecase | GetUserInfo | ошибка при decode cache")
			return nil, apperrors.ErrDB
		}
		return &cached, nil
	} else if err != redis.Nil {
		log.Printf("| usecase | GetUserInfo | ошибка при чтении cache")
		return nil, apperrors.ErrDB
	}

	account := &entities.Auth{PhoneNumber: dtoReq.PhoneNumber}
	userEntity, err := uc.trRepo.Get(ctx, account)
	if err != nil {
		return nil, err
	}

	isActive := strconv.FormatBool(userEntity.IsActive)
	res := dto.GetUserInfoResponse{
		ID:          userEntity.Id,
		PhoneNumber: dtoReq.PhoneNumber,
		FirstName:   userEntity.FirstName,
		LastName:    userEntity.LastName,
		RoleName:    userEntity.Role,
		IsActive:    &isActive,
		CreatedAt:   userEntity.CreatedAt,
		UpdatedAt:   userEntity.UpdatedAt,
	}

	raw, err := utils.MarshalToString(res)
	if err != nil {
		log.Printf("| usecase | GetUserInfo | ошибка при encode cache")
		return nil, apperrors.ErrDB
	}

	ttl := time.Duration(uc.cfg.JWT.AccessTTLMinutes) * time.Minute
	if err := uc.cacheRepo.Set(ctx, cacheKey, raw, ttl); err != nil {
		log.Printf("| usecase | GetUserInfo | ошибка при записе cache")
		return nil, apperrors.ErrDB
	}

	return &res, nil
}

func (uc *userUsecase) UpdateUserInfo(ctx context.Context, dtoReq dto.UpdateUserInfoRequest) (dtoRes *dto.UpdateUserInfoResponse, err error) {
	return nil, nil
}

func (uc *userUsecase) DeleteUserInfo(ctx context.Context, dtoReq dto.DeleteUserRequest) (dtoRes *dto.DeleteUserResponse, err error) {
	return nil, nil
}
