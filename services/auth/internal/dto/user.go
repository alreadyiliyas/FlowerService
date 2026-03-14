package dto

import "time"

type UserResponse struct {
	ID          *uint64    `json:"id,omitempty"`
	FirstName   *string    `json:"first_name,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type GetUserInfoRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
}

type GetUserInfoResponse struct {
	ID          *uint64    `json:"id,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	RoleName    *string    `json:"role_name,omitempty"`
	IsActive    *string    `json:"is_active,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type UpdateUserInfoRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	NewPhoneNumber *string `json:"new_phone_number,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

type UpdateUserInfoResponse struct {
	ID          *uint64    `json:"id,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	RoleName    *string    `json:"role_name,omitempty"`
	IsActive    *string    `json:"is_active,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type DeleteUserRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
}

type DeleteUserResponse struct {
	ID          *uint64    `json:"id,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	RoleName    *string    `json:"role_name,omitempty"`
	IsActive    *string    `json:"is_active,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
