package handlers

import (
	"errors"
	"net/http"

	"github.com/ilyas/flower/services/auth/internal/apperrors"
	"github.com/ilyas/flower/services/auth/internal/dto"
	"github.com/ilyas/flower/services/auth/internal/httpserver/middleware"
	userusecase "github.com/ilyas/flower/services/auth/internal/usecase/user"
	"github.com/ilyas/flower/services/auth/internal/utils"
)

type UserHandler struct {
	usecase userusecase.UsecaseUser
}

func NewUserHandler(uc userusecase.UsecaseUser) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	phone, ok := middleware.PhoneFromContext(r.Context())
	if !ok || phone == "" {
		utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
		return
	}

	req := dto.GetUserInfoRequest{PhoneNumber: &phone}
	resp, err := h.usecase.GetUserInfo(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidInput):
			utils.Send(w, http.StatusBadRequest, nil, err.Error())
		case errors.Is(err, apperrors.ErrNotFound), errors.Is(err, apperrors.ErrAccountNotFound), errors.Is(err, apperrors.ErrUserNotFound):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
		default:
			utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		}
		return
	}

	utils.Send(w, http.StatusOK, resp, "ok")
}
