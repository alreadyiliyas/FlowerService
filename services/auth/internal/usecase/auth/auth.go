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

	if err := ac.SendConfirmationCode(ctx, &dtoReq.PhoneNumber); err != nil {
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
	if !utils.IsValidPhoneNumber(*phone) {
		return fmt.Errorf("%w: проверьте номер телефона", apperrors.ErrInvalidInput)
	}

	code, err := utils.RandomConfirmCode()
	if err != nil {
		return fmt.Errorf("| usecase | SendConfirmationCode | %w", err)
	}

	if err := ac.cacheRepo.Set(ctx, utils.BuildConfirmKey(*phone), code, 15*time.Minute); err != nil {
		return fmt.Errorf("| usecase | SaveConfirmationCode | %w", err)
	}

	return nil
}

func (ac *authUsecase) VerifyAccount(ctx context.Context, dtoReq dto.VerifyAccountRequest) error {
	if dtoReq.PhoneNumber == nil || *dtoReq.PhoneNumber == "" {
		return fmt.Errorf("%w: задан пустой номер телефона", apperrors.ErrInvalidInput)
	}
	if dtoReq.Code == nil || *dtoReq.Code == "" {
		return fmt.Errorf("%w: задан пустой код подтверждения", apperrors.ErrInvalidInput)
	}
	if !utils.IsValidPhoneNumber(*dtoReq.PhoneNumber) {
		return fmt.Errorf("%w: проверьте номер телефона", apperrors.ErrInvalidInput)
	}

	key := utils.BuildConfirmKey(*dtoReq.PhoneNumber)
	cacheCode, err := ac.cacheRepo.Get(ctx, key)
	if err != nil {
		log.Printf("| usecase | VerifyAccount | read cache error: %v", err)
		return fmt.Errorf("%w: код подтверждения истек", apperrors.ErrNotFound)
	}

	if *dtoReq.Code != cacheCode {
		return fmt.Errorf("%w: неверный код подтверждения", apperrors.ErrInvalidInput)
	}

	if err := ac.trRepo.VerifyAccount(ctx, dtoReq.PhoneNumber); err != nil {
		log.Printf("| usecase | VerifyAccount | ошибка авторизации: %v", err)
		return fmt.Errorf("%w: %v", apperrors.ErrDB, err)
	}

	if err := ac.cacheRepo.Del(ctx, key); err != nil {
		log.Printf("| usecase | VerifyAccount | ошибка при удалении в кеше: %v", err)
	}

	return nil
}

func (ac *authUsecase) SetPassword(ctx context.Context, dtoReq dto.SetPasswordRequest) error {
	if dtoReq.PhoneNumber == nil || *dtoReq.PhoneNumber == "" {
		log.Printf("| usecase | SetPassword | задан пустой номер телефона")
		return fmt.Errorf("%w: задан пустой номер телефона", apperrors.ErrInvalidInput)
	}
	if dtoReq.Password == nil || *dtoReq.Password == "" {
		log.Printf("| usecase | SetPassword | проверка на валидацию пароля")
		return fmt.Errorf("%w: задан пустой пароль", apperrors.ErrInvalidInput)
	}
	if !utils.IsValidPhoneNumber(*dtoReq.PhoneNumber) {
		log.Printf("| usecase | SetPassword | проверка на валидацию номера телефона")
		return fmt.Errorf("%w: проверьте номер телефона", apperrors.ErrInvalidInput)
	}
	if err := utils.IsValidatePassword(*dtoReq.Password); err != nil {
		log.Printf("| usecase | SetPassword | проверка на валидацию пароля: %v", err)
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, err)
	}

	hash, err := utils.HashPassword(*dtoReq.Password)
	if err != nil {
		log.Printf("| usecase | SetPassword | не удалось захешировать пароль: %v", err)
		return fmt.Errorf("%w: не удалось захешировать пароль", apperrors.ErrDB)
	}

	account := &entities.Auth{PhoneNumber: dtoReq.PhoneNumber, PasswordHash: &hash}
	if err := ac.trRepo.SetPassword(ctx, account); err != nil {
		log.Printf("| usecase | SetPassword | set password error: %v", err)
		return fmt.Errorf("%w: %v", apperrors.ErrDB, err)
	}

	return nil
}

func (ac *authUsecase) SendCodeToUpdatePassword(ctx context.Context, phoneNumber *string) error {
	if phoneNumber == nil || *phoneNumber == "" {
		log.Printf("| usecase | SendCodeToUpdatePassword | задан пустой номер телефона")
		return fmt.Errorf("%w: задан пустой номер телефона", apperrors.ErrInvalidInput)
	}
	if !utils.IsValidPhoneNumber(*phoneNumber) {
		log.Printf("| usecase | SendCodeToUpdatePassword | проверка на валидацию номера телефона")
		return fmt.Errorf("%w: проверьте номер телефона", apperrors.ErrInvalidInput)
	}

	code, err := utils.RandomConfirmCode()
	if err != nil {
		log.Printf("| usecase | SendCodeToUpdatePassword | error when generate code: %w", err)
		return fmt.Errorf("%w: %v", apperrors.ErrDB, err)
	}

	if err := ac.cacheRepo.Set(ctx, utils.BuildPasswordUpdateKey(*phoneNumber), code, 30*time.Minute); err != nil {
		log.Printf("| usecase | SendCodeToUpdatePassword | error to set code: %w", err)
		return fmt.Errorf("%w: %v", apperrors.ErrDB, err)
	}

	return nil
}

func (ac *authUsecase) UpdatePassword(ctx context.Context, dtoReq dto.ConfirmUpdatePasswordRequest) error {
	if dtoReq.PhoneNumber == nil || *dtoReq.PhoneNumber == "" {
		log.Printf("| usecase | UpdatePassword | задан пустой номер телефона")
		return fmt.Errorf("%w: задан пустой номер телефона", apperrors.ErrInvalidInput)
	}
	if dtoReq.Code == nil || *dtoReq.Code == "" {
		log.Printf("| usecase | UpdatePassword | задан пустой код подтверждения")
		return fmt.Errorf("%w: задан пустой код подтверждения", apperrors.ErrInvalidInput)
	}
	if dtoReq.NewPassword == nil || *dtoReq.NewPassword == "" {
		log.Printf("| usecase | UpdatePassword | задан пустой пароль")
		return fmt.Errorf("%w: задан пустой пароль", apperrors.ErrInvalidInput)
	}
	if !utils.IsValidPhoneNumber(*dtoReq.PhoneNumber) {
		log.Printf("| usecase | UpdatePassword | проверка на валидацию номера телефона")
		return fmt.Errorf("%w: проверьте номер телефона", apperrors.ErrInvalidInput)
	}
	if err := utils.IsValidatePassword(*dtoReq.NewPassword); err != nil {
		log.Printf("| usecase | UpdatePassword | проверка на валидацию пароля: %v", err)
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, err)
	}

	key := utils.BuildPasswordUpdateKey(*dtoReq.PhoneNumber)
	cacheCode, err := ac.cacheRepo.Get(ctx, key)
	if err != nil {
		log.Printf("| usecase | UpdatePassword | read cache error: %v", err)
		return fmt.Errorf("%w: код подтверждения истек", apperrors.ErrNotFound)
	}

	if *dtoReq.Code != cacheCode {
		return fmt.Errorf("%w: неверный код подтверждения", apperrors.ErrInvalidInput)
	}

	hash, err := utils.HashPassword(*dtoReq.NewPassword)
	if err != nil {
		return fmt.Errorf("%w: не удалось захешировать пароль", apperrors.ErrDB)
	}

	account := &entities.Auth{PhoneNumber: dtoReq.PhoneNumber, PasswordHash: &hash}
	if err := ac.trRepo.UpdatePassword(ctx, account); err != nil {
		log.Printf("| usecase | update password | db error: %v", err)
		return fmt.Errorf("%w: не установить пароль", apperrors.ErrDB)
	}

	if err := ac.cacheRepo.Del(ctx, key); err != nil {
		log.Printf("| usecase | update password | delete code error: %v", err)
	}

	return nil
}
