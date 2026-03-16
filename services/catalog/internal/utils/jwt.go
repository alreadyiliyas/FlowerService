package utils

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	UserID    uint64 `json:"user_id"`
	Role      string `json:"role,omitempty"`
	Phone     string `json:"phone,omitempty"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

func ParseAccessToken(tokenString, secret string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
