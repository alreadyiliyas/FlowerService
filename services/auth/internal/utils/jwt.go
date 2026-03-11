package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID uint64 `json:"user_id"`
	Role   string `json:"role,omitempty"`
	Phone  string `json:"phone,omitempty"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint64, role, phone, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		UserID: userID,
		Role:   role,
		Phone:  phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
