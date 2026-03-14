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
		return nil, fmt.Errorf("%w: phone number is empty", apperrors.ErrInvalidInput)
	}
	if !utils.IsValidPhoneNumber(*dtoReq.PhoneNumber) {
		return nil, fmt.Errorf("%w: invalid phone number", apperrors.ErrInvalidInput)
	}

	cacheKey := utils.BuildUserInfoKey(*dtoReq.PhoneNumber)
	if raw, err := uc.cacheRepo.Get(ctx, cacheKey); err == nil {
		var cached dto.GetUserInfoResponse
		if err := utils.UnmarshalFromString(raw, &cached); err != nil {
			log.Printf("| usecase | GetUserInfo | decode cache error: %v", err)
			return nil, apperrors.ErrDB
		}
		return &cached, nil
	} else if err != redis.Nil {
		log.Printf("| usecase | GetUserInfo | read cache error: %v", err)
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
		PhoneNumber: userEntity.PhoneNumber,
		FirstName:   userEntity.FirstName,
		LastName:    userEntity.LastName,
		RoleName:    userEntity.Role,
		IsActive:    &isActive,
		AvatarURL:   userEntity.AvatarURL,
		CreatedAt:   userEntity.CreatedAt,
		UpdatedAt:   userEntity.UpdatedAt,
	}

	raw, err := utils.MarshalToString(res)
	if err != nil {
		log.Printf("| usecase | GetUserInfo | encode cache error: %v", err)
		return nil, apperrors.ErrDB
	}

	ttl := time.Duration(uc.cfg.JWT.AccessTTLMinutes) * time.Minute
	if err := uc.cacheRepo.Set(ctx, cacheKey, raw, ttl); err != nil {
		log.Printf("| usecase | GetUserInfo | write cache error: %v", err)
		return nil, apperrors.ErrDB
	}

	return &res, nil
}

func (uc *userUsecase) UpdateUserInfo(ctx context.Context, dtoReq dto.UpdateUserInfoRequest) (dtoRes *dto.UpdateUserInfoResponse, err error) {
	if dtoReq.PhoneNumber == nil || *dtoReq.PhoneNumber == "" {
		return nil, fmt.Errorf("%w: phone number is empty", apperrors.ErrInvalidInput)
	}
	if !utils.IsValidPhoneNumber(*dtoReq.PhoneNumber) {
		return nil, fmt.Errorf("%w: invalid phone number", apperrors.ErrInvalidInput)
	}
	if dtoReq.NewPhoneNumber != nil && *dtoReq.NewPhoneNumber != "" && !utils.IsValidPhoneNumber(*dtoReq.NewPhoneNumber) {
		return nil, fmt.Errorf("%w: invalid new phone number", apperrors.ErrInvalidInput)
	}
	if dtoReq.FirstName == nil && dtoReq.LastName == nil && dtoReq.NewPhoneNumber == nil && dtoReq.AvatarURL == nil {
		return nil, fmt.Errorf("%w: nothing to update", apperrors.ErrInvalidInput)
	}

	// Load current values to avoid overwriting with nil.
	current, err := uc.trRepo.Get(ctx, &entities.Auth{PhoneNumber: dtoReq.PhoneNumber})
	if err != nil {
		return nil, err
	}

	firstName := current.FirstName
	lastName := current.LastName
	avatarURL := current.AvatarURL
	newPhone := current.PhoneNumber

	if dtoReq.FirstName != nil {
		firstName = dtoReq.FirstName
	}
	if dtoReq.LastName != nil {
		lastName = dtoReq.LastName
	}
	if dtoReq.AvatarURL != nil {
		avatarURL = dtoReq.AvatarURL
	}
	if dtoReq.NewPhoneNumber != nil && *dtoReq.NewPhoneNumber != "" {
		newPhone = dtoReq.NewPhoneNumber
	}

	account := &entities.Auth{
		PhoneNumber: dtoReq.PhoneNumber,
		User: entities.User{
			FirstName:   firstName,
			LastName:    lastName,
			AvatarURL:   avatarURL,
			PhoneNumber: newPhone,
		},
	}

	userEntity, err := uc.trRepo.Update(ctx, account)
	if err != nil {
		return nil, err
	}

	isActive := strconv.FormatBool(userEntity.IsActive)
	res := dto.UpdateUserInfoResponse{
		ID:          userEntity.Id,
		PhoneNumber: userEntity.PhoneNumber,
		FirstName:   userEntity.FirstName,
		LastName:    userEntity.LastName,
		RoleName:    userEntity.Role,
		IsActive:    &isActive,
		AvatarURL:   userEntity.AvatarURL,
		UpdatedAt:   userEntity.UpdatedAt,
	}

	// Update user info cache (key may change).
	oldCacheKey := utils.BuildUserInfoKey(*dtoReq.PhoneNumber)
	newCacheKey := utils.BuildUserInfoKey(*userEntity.PhoneNumber)
	cachePayload := dto.GetUserInfoResponse{
		ID:          res.ID,
		PhoneNumber: res.PhoneNumber,
		FirstName:   res.FirstName,
		LastName:    res.LastName,
		RoleName:    res.RoleName,
		IsActive:    res.IsActive,
		AvatarURL:   res.AvatarURL,
		CreatedAt:   userEntity.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}

	raw, err := utils.MarshalToString(cachePayload)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode cache", apperrors.ErrDB)
	}

	ttl := time.Duration(uc.cfg.JWT.AccessTTLMinutes) * time.Minute
	if err := uc.cacheRepo.Set(ctx, newCacheKey, raw, ttl); err != nil {
		return nil, fmt.Errorf("%w: failed to write cache", apperrors.ErrDB)
	}
	if oldCacheKey != newCacheKey {
		_ = uc.cacheRepo.Del(ctx, oldCacheKey)
	}

	// If phone changed, migrate session index and cached session payloads.
	if dtoReq.NewPhoneNumber != nil && *dtoReq.NewPhoneNumber != "" && *dtoReq.NewPhoneNumber != *dtoReq.PhoneNumber {
		oldSessionKey := utils.BuildSessionKeyByPhone(*dtoReq.PhoneNumber)
		newSessionKey := utils.BuildSessionKeyByPhone(*dtoReq.NewPhoneNumber)
		sessionIDs, err := uc.cacheRepo.SMembers(ctx, oldSessionKey)
		if err == nil && len(sessionIDs) > 0 {
			for _, sid := range sessionIDs {
				if sid == "" {
					continue
				}
				sKey := utils.BuildSessionKey(sid)
				rawSession, err := uc.cacheRepo.Get(ctx, sKey)
				if err != nil {
					continue
				}
				var session dto.RefreshCache
				if err := utils.UnmarshalFromString(rawSession, &session); err != nil {
					continue
				}
				session.PhoneNumber = *dtoReq.NewPhoneNumber
				if firstName != nil {
					session.FirstName = *firstName
				}
				if lastName != nil {
					session.LastName = *lastName
				}
				sessionRaw, err := utils.MarshalToString(session)
				if err != nil {
					continue
				}
				_ = uc.cacheRepo.Set(ctx, sKey, sessionRaw, time.Duration(uc.cfg.JWT.RefreshTTLDays)*24*time.Hour)
				_ = uc.cacheRepo.SAdd(ctx, newSessionKey, sid)
			}
			_ = uc.cacheRepo.Del(ctx, oldSessionKey)
		}
		_ = uc.cacheRepo.Expire(ctx, newSessionKey, time.Duration(uc.cfg.JWT.RefreshTTLDays)*24*time.Hour)
	}

	return &res, nil
}

func (uc *userUsecase) DeleteUserInfo(ctx context.Context, dtoReq dto.DeleteUserRequest) (dtoRes *dto.DeleteUserResponse, err error) {
	return nil, nil
}
