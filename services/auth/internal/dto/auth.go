package dto

import "time"

type RegistrationRequest struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Role        string `json:"role,omitempty"`
}

type RegistrationResponse struct {
	UserID      *uint64    `json:"user_id,omitempty"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	Role        *string    `json:"role,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

type SetPasswordRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
	Password    *string `json:"password,omitempty"`
}

type UpdatePasswordRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
	OldPassword *string `json:"old_password,omitempty"`
	NewPassword *string `json:"new_password,omitempty"`
}

type VerifyAccountRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
	Code        *string `json:"code,omitempty"`
}

type LoginRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
	Password    *string `json:"password,omitempty"`
}

type LoginResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshTokenKey string `json:"refresh_token_key"`
}

type RefreshTokenRequest struct {
	RefreshTokenKey *string `json:"refresh_token_key"`
}

type RefreshTokenResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshTokenKey string `json:"refresh_token_key"`
}

type LogoutRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
}
