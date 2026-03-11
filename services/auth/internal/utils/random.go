package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

func RandomConfirmCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("%w: %v", err, "ошибка при генерировании рандомного числа")
	}

	return fmt.Sprintf("%06d", n.Uint64()), nil
}

func RandomToken(bytesLen int) (string, error) {
	if bytesLen <= 0 {
		return "", fmt.Errorf("bytesLen must be > 0")
	}
	b := make([]byte, bytesLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
