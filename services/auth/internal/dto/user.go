package dto

import "time"

type UserResponse struct {
	ID          *uint64    `json:"id,omitempty"`
	FirstName   *string    `json:"first_name,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
