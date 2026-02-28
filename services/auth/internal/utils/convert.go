package utils

import (
	"fmt"
	"math"
)

func ToString(v interface{}, field string) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("invalid %s type: got %T, want string", field, v)
	}
	return s, nil
}

func ToUint64(v interface{}, field string) (uint64, error) {
	switch n := v.(type) {
	case uint64:
		return n, nil
	case uint32:
		return uint64(n), nil
	case uint16:
		return uint64(n), nil
	case uint8:
		return uint64(n), nil
	case uint:
		return uint64(n), nil
	case int64:
		if n < 0 {
			return 0, fmt.Errorf("invalid %s value: got %d, must be >= 0", field, n)
		}
		return uint64(n), nil
	case int32:
		if n < 0 {
			return 0, fmt.Errorf("invalid %s value: got %d, must be >= 0", field, n)
		}
		return uint64(n), nil
	case int16:
		if n < 0 {
			return 0, fmt.Errorf("invalid %s value: got %d, must be >= 0", field, n)
		}
		return uint64(n), nil
	case int8:
		if n < 0 {
			return 0, fmt.Errorf("invalid %s value: got %d, must be >= 0", field, n)
		}
		return uint64(n), nil
	case int:
		if n < 0 {
			return 0, fmt.Errorf("invalid %s value: got %d, must be >= 0", field, n)
		}
		return uint64(n), nil
	case float64:
		if n < 0 || math.Trunc(n) != n {
			return 0, fmt.Errorf("invalid %s value: got %v, must be non-negative integer", field, n)
		}
		return uint64(n), nil
	default:
		return 0, fmt.Errorf("invalid %s type: got %T, want unsigned integer", field, v)
	}
}
