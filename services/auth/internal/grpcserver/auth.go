package grpcserver

import (
	"context"
	"errors"

	grpcauth "github.com/ilyas/flower/pkg/grpc/authcontext"
	"github.com/ilyas/flower/services/auth/internal/apperrors"
	"github.com/ilyas/flower/services/auth/internal/dto"
	authusecase "github.com/ilyas/flower/services/auth/internal/usecase/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct {
	authUC authusecase.UsecaseAuth
}

// NewAuthHandler адаптирует usecase auth к gRPC интерфейсу.
func NewAuthHandler(authUC authusecase.UsecaseAuth) grpcauth.AuthServiceServer {
	return &authHandler{authUC: authUC}
}

func (h *authHandler) GetUserContext(ctx context.Context, req *grpcauth.GetUserContextRequest) (*grpcauth.GetUserContextResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	resp, err := h.authUC.GetUserContext(ctx, dto.GetUserContextRequest{
		AccessToken: req.AccessToken,
		SessionID:   req.SessionID,
	})
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidInput):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, apperrors.ErrUnauthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &grpcauth.GetUserContextResponse{
		UserID:          resp.UserID,
		Role:            resp.Role,
		PhoneNumber:     resp.PhoneNumber,
		SessionID:       resp.SessionID,
		FirstName:       resp.FirstName,
		LastName:        resp.LastName,
		IsAuthenticated: resp.IsAuthenticated,
	}, nil
}
