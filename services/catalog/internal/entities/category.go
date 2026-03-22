package entities

import "time"

type Category struct {
	ID          *uint64    `json:"id"`
	Name        *string    `json:"name"`
	Slug        *string    `json:"slug"`
	Description *string    `json:"description,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
