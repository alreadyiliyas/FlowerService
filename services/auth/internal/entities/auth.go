package entities

import "time"

type Auth struct {
	Id           *uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       *uint64 `json:"user_id"`
	PhoneNumber  *string `json:"phone_number"`
	PasswordHash *string `json:"-"`
	User         User    `gorm:"foreignKey:UserId"`
}

type Token struct {
	Id           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       uint64    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at"`
}
