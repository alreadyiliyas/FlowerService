package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ilyas/flower/services/auth/internal/apperrors"
	"github.com/ilyas/flower/services/auth/internal/dto"
	auth "github.com/ilyas/flower/services/auth/internal/usecase/auth"
	"github.com/ilyas/flower/services/auth/internal/utils"
)

// AuthHandler обрабатывает HTTP-запросы для аутентификации.
type AuthHandler struct {
	usecase auth.Usecase
}

func NewAuthHandler(uc auth.Usecase) *AuthHandler {
	return &AuthHandler{usecase: uc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.usecase.Registration(context.Background(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidInput):
			utils.Send(w, http.StatusBadRequest, nil, err.Error())
		case errors.Is(err, apperrors.ErrDuplicatePhone):
			utils.Send(w, http.StatusConflict, nil, err.Error())
		case errors.Is(err, apperrors.ErrNotFound):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
		default:
			utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		}
	}

	utils.Send(w, http.StatusCreated, resp, "пользователь успешно авторизован")
}

func (h *AuthHandler) SetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.SetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: валидация (phone_number и password обязательны)

	// TODO: проверка, что пользователь существует

	// TODO: хеширование пароля (bcrypt)

	// TODO: сохранение password_hash в authRepo

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "password set"})
}

func (h *AuthHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req dto.SetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: валидация (phone_number и password обязательны)

	// TODO: проверка, что пользователь существует

	// TODO: хеширование пароля (bcrypt)

	// TODO: сохранение password_hash в authRepo

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "password set"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: валидация (phone_number и password обязательны)

	// TODO: найти пользователя по phone_number через authRepo

	// TODO: сравнить password_hash с введённым паролем (bcrypt.CompareHashAndPassword)

	// TODO: если пароль верный:
	//   - сгенерировать access_token (JWT)
	//   - создать refresh_token и сохранить в БД
	//   - вернуть LoginResponse

	// TODO: если пароль неверный - вернуть 401 Unauthorized

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) VerifyAccount(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй верификацию по SMS-коду
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуй обновление access_token по refresh_token_key
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

}
