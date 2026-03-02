package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ilyas/flower/services/auth/internal/apperrors"
	"github.com/ilyas/flower/services/auth/internal/config"
	"github.com/ilyas/flower/services/auth/internal/dto"
	"github.com/ilyas/flower/services/auth/internal/entities"
	auth "github.com/ilyas/flower/services/auth/internal/repositories/auth"
	cache "github.com/ilyas/flower/services/auth/internal/repositories/cache"
	"github.com/ilyas/flower/services/auth/internal/utils"
)

type authUsecase struct {
	cfg       *config.Config
	trRepo    auth.AuthRepository
	cacheRepo cache.CacheRepository
}

func New(cfg *config.Config, ar auth.AuthRepository, cr cache.CacheRepository) Usecase {
	return &authUsecase{
		cfg:       cfg,
		trRepo:    ar,
		cacheRepo: cr,
	}
}

func (ac *authUsecase) Registration(ctx context.Context, dtoReq dto.RegistrationRequest) (dtoRes *dto.RegistrationResponse, err error) {
	if err := utils.ValidateAuthData(dtoReq); err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, err.Error())
	}
	if dtoReq.Role == "" {
		dtoReq.Role = "user"
	}
	userEntity := entities.User{
		FirstName: &dtoReq.FirstName,
		LastName:  &dtoReq.LastName,
		Role:      &dtoReq.Role,
		IsActive:  false,
		Version:   1,
	}

	accEntity := entities.Auth{
		PhoneNumber: &dtoReq.PhoneNumber,
	}

	createdUser, err := ac.trRepo.CreateUser(ctx, &userEntity, &accEntity)
	if err != nil {
		log.Printf("| usecase | create user error: %v", err)
		return nil, err
	}

	err = ac.SendConfirmationCode(ctx, &dtoReq.PhoneNumber)
	if err != nil {
		return nil, err
	}

	res := dto.RegistrationResponse{
		UserID:      createdUser.Id,
		FirstName:   createdUser.FirstName,
		LastName:    createdUser.LastName,
		Role:        createdUser.Role,
		PhoneNumber: accEntity.PhoneNumber,
		CreatedAt:   createdUser.CreatedAt,
	}

	return &res, nil
}

func (ac *authUsecase) SendConfirmationCode(ctx context.Context, phone *string) error {
	if phone == nil || *phone == "" {
		return fmt.Errorf("%w: нельзя задать пустой номер телефона", apperrors.ErrInvalidInput)
	}

	code, err := utils.RandomConfirmCode()
	if err != nil {
		return fmt.Errorf("| usecase | SendConfirmationCode | %w", err)
	}

	ttl := time.Duration(15 * time.Minute)

	err = ac.cacheRepo.SaveConfirmationCode(ctx, phone, &code, ttl)
	if err != nil {
		return fmt.Errorf("| usecase | SaveConfirmationCode | %w", err)
	}

	return nil
}

func (ac *authUsecase) VerifyAccount(ctx context.Context, dtoReq dto.VerifyAccountRequest) error {
	if dtoReq.PhoneNumber == nil || *dtoReq.PhoneNumber == "" {
		return fmt.Errorf("%w: задан пустой номер телефона", apperrors.ErrInvalidInput)
	}

	cacheCode, err := ac.cacheRepo.GetConfirmCode(ctx, dtoReq.PhoneNumber)
	if err != nil {
		log.Printf("| usecase | verify account | read cache error: %v", err)
		return fmt.Errorf("%w: код подтверждения истек", apperrors.ErrNotFound)
	}

	if *dtoReq.Code != cacheCode {
		return fmt.Errorf("%w: не верный код подтверждения", apperrors.ErrInvalidInput)
	}

	err = ac.trRepo.VerifyAccount(ctx, dtoReq.PhoneNumber)
	if err != nil {
		log.Printf("| usecase | verify account | activate in tarantoolDB error: %v", err)
		return err
	}

	if err := ac.cacheRepo.DeleteConfirmCode(ctx, dtoReq.PhoneNumber); err != nil {
		log.Printf("| usecase | delete confirm code | error: %v", err)
	}

	return nil
}

func (ac *authUsecase) SetPassword(ctx context.Context, dtoReq dto.SetPasswordRequest) error {
	if dtoReq.PhoneNumber == nil || *dtoReq.PhoneNumber == "" {
		return fmt.Errorf("%w: задан пустой номер телефона", apperrors.ErrInvalidInput)
	}

	err := utils.IsValidatePassword(*dtoReq.Password)
	if err != nil {
		return err
	}

	hash, err := utils.HashPassword(*dtoReq.Password)
	if err != nil {
		return fmt.Errorf("%w: не удалось захешировать пароль", apperrors.ErrDB)
	}
	account := &entities.Auth{
		PhoneNumber:  dtoReq.PhoneNumber,
		PasswordHash: &hash,
	}

	if err = ac.trRepo.SetPassword(ctx, account); err != nil {
		log.Printf("| usecase | set password | set password error: %v", err)
		return err
	}

	return nil
}
