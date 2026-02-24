package entities

import "time"

type User struct {
	Id        *uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	Role      *string    `json:"role"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	IsActive  bool       `json:"is_active"`
	Version   int        `gorm:"default:0"`
}
