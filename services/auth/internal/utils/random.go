package utils

import (
	"crypto/rand"
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
