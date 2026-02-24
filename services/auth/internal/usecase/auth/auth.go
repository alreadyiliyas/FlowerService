package usecase

import (
	"context"
	"fmt"

	"github.com/ilyas/flower/services/auth/internal/apperrors"
	"github.com/ilyas/flower/services/auth/internal/config"
	"github.com/ilyas/flower/services/auth/internal/dto"
	"github.com/ilyas/flower/services/auth/internal/entities"
	auth "github.com/ilyas/flower/services/auth/internal/repositories/auth"
	"github.com/ilyas/flower/services/auth/internal/utils"
)

type authUsecase struct {
	cfg    *config.Config
	trRepo auth.AuthRepository
}

func New(cfg *config.Config, repo auth.AuthRepository) Usecase {
	return &authUsecase{
		cfg:    cfg,
		trRepo: repo,
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
		return nil, err
	}

	// err = ac.SendConfirmationCode(ctx, &dtoReq.PhoneNumber)
	// if err != nil {
	// 	return nil, err
	// }

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
