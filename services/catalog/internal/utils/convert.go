package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ilyas/flower/services/catalog/internal/dto"
	"github.com/ilyas/flower/services/catalog/internal/entities"
)

func ToString(v interface{}, field string) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("invalid %s type: got %T, want string", field, v)
	}
	return s, nil
}

func ToStringNullable(v interface{}, field string) (*string, error) {
	if v == nil {
		return nil, nil
	}
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("invalid %s type: got %T, want string", field, v)
	}
	return &s, nil
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

func ToInt(v interface{}, field string) (int, error) {
	switch n := v.(type) {
	case int64:
		return int(n), nil
	case int32:
		return int(n), nil
	case int16:
		return int(n), nil
	case int8:
		return int(n), nil
	case int:
		return int(n), nil
	case string:
		atoi, err := strconv.Atoi(n)
		if err != nil {
			return 0, fmt.Errorf("%s: can't convert to int from %T type", field, v)
		}
		return atoi, nil
	default:
		return 0, fmt.Errorf("invalid %s type: got %T, want unsigned integer", field, v)
	}
}

func ToBool(v interface{}, field string) (bool, error) {
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("invalid %s type: got %T, want string", field, v)
	}
	return b, nil
}

func ValueOrNull(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func MarshalToString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func UnmarshalFromString(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

func MapCategoryToDTO(item entities.Category) dto.Category {
	category := dto.Category{}
	if item.ID != nil {
		category.ID = *item.ID
	}
	if item.Name != nil {
		category.Name = *item.Name
	}
	if item.Slug != nil {
		category.Slug = *item.Slug
	}
	if item.Description != nil {
		category.Description = *item.Description
	}
	if item.ImageURL != nil {
		category.ImageURL = *item.ImageURL
	}
	if item.CreatedAt != nil {
		category.CreatedAt = item.CreatedAt.Format(time.RFC3339)
	}
	if item.UpdatedAt != nil {
		category.UpdatedAt = item.UpdatedAt.Format(time.RFC3339)
	}
	return category
}
