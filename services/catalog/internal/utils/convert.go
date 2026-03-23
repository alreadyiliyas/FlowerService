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
	case uint64:
		return int(n), nil
	case uint32:
		return int(n), nil
	case uint16:
		return int(n), nil
	case uint8:
		return int(n), nil
	case uint:
		return int(n), nil
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

func Uint64OrNull(v *uint64) interface{} {
	if v == nil {
		return nil
	}
	return *v
}

func BoolOrNull(v *bool) interface{} {
	if v == nil {
		return nil
	}
	return *v
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

func MapProductEntityToDTO(item entities.Product) *dto.Product {
	product := &dto.Product{
		IsAvailable: item.IsAvailable,
		Images:      append([]string(nil), item.Images...),
	}

	if item.ID != nil {
		product.ID = *item.ID
	}
	if item.Name != nil {
		product.Name = *item.Name
	}
	if item.Description != nil {
		product.Description = *item.Description
	}
	if item.CategoryID != nil {
		product.CategoryID = *item.CategoryID
	}
	if item.SellerID != nil {
		product.SellerID = *item.SellerID
	}
	if item.MainImageURL != nil {
		product.MainImageURL = *item.MainImageURL
	}
	if item.Currency != nil {
		product.Currency = *item.Currency
	}
	if item.PricePerStem != nil {
		product.PricePerStem = int(*item.PricePerStem)
	}
	if item.MinStems != nil {
		product.MinStems = int(*item.MinStems)
	}
	if item.MaxStems != nil {
		product.MaxStems = int(*item.MaxStems)
	}
	if item.CreatedAt != nil {
		product.CreatedAt = item.CreatedAt.Format(time.RFC3339)
	}
	if item.UpdatedAt != nil {
		product.UpdatedAt = item.UpdatedAt.Format(time.RFC3339)
	}

	product.Sizes = make([]dto.SizePrice, 0, len(item.Sizes))
	for _, size := range item.Sizes {
		sizeDTO := dto.SizePrice{}
		if size.Size != nil {
			sizeDTO.Size = *size.Size
		}
		if size.BasePrice != nil {
			sizeDTO.BasePrice = int(*size.BasePrice)
		}
		product.Sizes = append(product.Sizes, sizeDTO)
	}

	product.Composition = make([]dto.CompositionItem, 0, len(item.Composition))
	for _, composition := range item.Composition {
		compositionDTO := dto.CompositionItem{}
		if composition.FlowerType != nil {
			compositionDTO.FlowerType = *composition.FlowerType
		}
		if composition.Stems != nil {
			compositionDTO.Stems = int(*composition.Stems)
		}
		product.Composition = append(product.Composition, compositionDTO)
	}

	if item.Discount != nil {
		product.Discount = &dto.Discount{}
		if item.Discount.Type != nil {
			product.Discount.Type = *item.Discount.Type
		}
		if item.Discount.Value != nil {
			product.Discount.Value = int(*item.Discount.Value)
		}
		if item.Discount.StartsAt != nil {
			product.Discount.StartsAt = *item.Discount.StartsAt
		}
		if item.Discount.EndsAt != nil {
			product.Discount.EndsAt = *item.Discount.EndsAt
		}
	}

	return product
}

func MapProductFilterToEntity(filter dto.ProductFilter) entities.ProductFilter {
	entityFilter := entities.ProductFilter{
		Page:     1,
		PageSize: 20,
	}

	if filter.CategoryID != nil {
		entityFilter.CategoryID = filter.CategoryID
	}
	if filter.SellerID != nil {
		entityFilter.SellerID = filter.SellerID
	}
	if filter.PriceMin != nil {
		value := uint64(*filter.PriceMin)
		entityFilter.PriceMin = &value
	}
	if filter.PriceMax != nil {
		value := uint64(*filter.PriceMax)
		entityFilter.PriceMax = &value
	}
	if filter.Size != nil {
		size := *filter.Size
		entityFilter.Size = &size
	}
	if filter.IsAvailable != nil {
		entityFilter.IsAvailable = filter.IsAvailable
	}
	if filter.Page != nil && *filter.Page > 0 {
		entityFilter.Page = *filter.Page
	}
	if filter.PageSize != nil && *filter.PageSize > 0 {
		entityFilter.PageSize = *filter.PageSize
	}

	return entityFilter
}

func MapPaginatedProductsToDTO(page entities.PaginatedProducts) dto.PaginatedProducts {
	result := dto.PaginatedProducts{
		Items:    make([]dto.Product, 0, len(page.Items)),
		Total:    page.Total,
		Page:     page.Page,
		PageSize: page.PageSize,
	}

	for _, item := range page.Items {
		result.Items = append(result.Items, *MapProductEntityToDTO(item))
	}

	return result
}
