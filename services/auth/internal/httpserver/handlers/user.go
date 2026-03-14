package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

func (h *UserHandler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	phone, ok := middleware.PhoneFromContext(r.Context())
	if !ok || phone == "" {
		utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
		return
	}

	var req dto.UpdateUserInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.PhoneNumber = &phone

	resp, err := h.usecase.UpdateUserInfo(r.Context(), req)
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

	utils.Send(w, http.StatusOK, resp, "пользователь обновлен")
}

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	phone, ok := middleware.PhoneFromContext(r.Context())
	if !ok || phone == "" {
		utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
		return
	}

	const maxUploadSize = 5 << 20 // 5MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		utils.Send(w, http.StatusBadRequest, nil, "invalid multipart form")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		utils.Send(w, http.StatusBadRequest, nil, "avatar file is required")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		utils.Send(w, http.StatusBadRequest, nil, "unsupported file type")
		return
	}

	if err := os.MkdirAll("public/avatars", 0o755); err != nil {
		utils.Send(w, http.StatusInternalServerError, nil, "failed to create directory")
		return
	}

	name, err := utils.RandomToken(16)
	if err != nil {
		utils.Send(w, http.StatusInternalServerError, nil, "failed to generate filename")
		return
	}

	fileName := name + ext
	dstPath := filepath.Join("public", "avatars", fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		utils.Send(w, http.StatusInternalServerError, nil, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.Send(w, http.StatusInternalServerError, nil, "failed to save file")
		return
	}

	avatarURL := "/public/avatars/" + fileName
	req := dto.UpdateUserInfoRequest{
		PhoneNumber: &phone,
		AvatarURL:   &avatarURL,
	}

	resp, err := h.usecase.UpdateUserInfo(r.Context(), req)
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

	utils.Send(w, http.StatusOK, resp, "Аватар обновлен")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	phone, ok := middleware.PhoneFromContext(r.Context())
	if !ok || phone == "" {
		utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
		return
	}

	var req dto.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.PhoneNumber = &phone

	err := h.usecase.DeleteUserInfo(r.Context(), req)
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

	utils.Send(w, http.StatusOK, nil, "пользователь удален")
}
