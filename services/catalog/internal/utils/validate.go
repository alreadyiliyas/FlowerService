package utils

import (
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/dto"
	"github.com/ilyas/flower/services/catalog/internal/entities"
)

var categorySlugRegexp = regexp.MustCompile(`^[a-z0-9-]+$`)

func ValidateCategory(in dto.CreateCategoryRequest) (*entities.Category, error) {
	name := strings.TrimSpace(in.Name)
	slug := strings.TrimSpace(strings.ToLower(in.Slug))
	description := strings.TrimSpace(in.Description)

	switch {
	case name == "":
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "задано пустое имя для категории")
	case slug == "":
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "задан пустой slug")
	case !categorySlugRegexp.MatchString(slug):
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "slug должен содержать только латиницу, цифры и дефис")
	case in.Image == nil || in.ImageHeader == nil:
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "передана пустая картинка")
	}

	return &entities.Category{
		Name:        &name,
		Slug:        &slug,
		Description: &description,
	}, nil
}

func ValidateCategoryUpdate(in dto.UpdateCategoryRequest) (*entities.Category, error) {
	name := strings.TrimSpace(in.Name)
	slug := strings.TrimSpace(strings.ToLower(in.Slug))
	description := strings.TrimSpace(in.Description)

	if slug != "" && !categorySlugRegexp.MatchString(slug) {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "slug должен содержать только латиницу, цифры и дефис")
	}

	return &entities.Category{
		Name:        stringPtrOrNil(name),
		Slug:        stringPtrOrNil(slug),
		Description: stringPtrOrNil(description),
	}, nil
}

func ValidateProduct(in dto.Product) (*entities.Product, error) {
	name := strings.TrimSpace(in.Name)
	description := strings.TrimSpace(in.Description)
	currency := strings.TrimSpace(strings.ToUpper(in.Currency))

	switch {
	case name == "":
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "задано пустое имя продукта")
	case in.CategoryID == 0:
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "не передан category_id")
	case in.SellerID == 0:
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "не передан seller_id")
	case currency == "":
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "не передана валюта")
	case len(in.Sizes) == 0:
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "не передан список sizes")
	case in.PricePerStem <= 0:
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "price_per_stem должен быть больше 0")
	case in.MinStems <= 0:
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "min_stems должен быть больше 0")
	case in.MaxStems <= 0:
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "max_stems должен быть больше 0")
	case in.MinStems > in.MaxStems:
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "min_stems не может быть больше max_stems")
	case len(in.Composition) == 0:
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "не передан composition")
	}

	sizes := make([]entities.SizePrice, 0, len(in.Sizes))
	for _, item := range in.Sizes {
		size := strings.TrimSpace(strings.ToUpper(item.Size))
		switch size {
		case "S", "M", "L":
		default:
			return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "size должен быть одним из: S, M, L")
		}
		if item.BasePrice <= 0 {
			return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "base_price должен быть больше 0")
		}

		basePrice := uint64(item.BasePrice)
		sizes = append(sizes, entities.SizePrice{
			Size:      &size,
			BasePrice: &basePrice,
		})
	}

	composition := make([]entities.CompositionItem, 0, len(in.Composition))
	for _, item := range in.Composition {
		flowerType := strings.TrimSpace(item.FlowerType)
		if flowerType == "" {
			return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "flower_type не может быть пустым")
		}
		if item.Stems <= 0 {
			return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "stems должен быть больше 0")
		}

		stems := uint64(item.Stems)
		composition = append(composition, entities.CompositionItem{
			FlowerType: &flowerType,
			Stems:      &stems,
		})
	}

	product := &entities.Product{
		Name:         &name,
		Description:  stringPtrOrNil(description),
		CategoryID:   uint64Ptr(in.CategoryID),
		SellerID:     uint64Ptr(in.SellerID),
		IsAvailable:  in.IsAvailable,
		Currency:     &currency,
		Sizes:        sizes,
		PricePerStem: uint64Ptr(uint64(in.PricePerStem)),
		MinStems:     uint64Ptr(uint64(in.MinStems)),
		MaxStems:     uint64Ptr(uint64(in.MaxStems)),
		Composition:  composition,
	}

	if in.Discount != nil {
		discountType := strings.TrimSpace(in.Discount.Type)
		if discountType == "" {
			return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "discount.type не может быть пустым")
		}
		if in.Discount.Value <= 0 {
			return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "discount.value должен быть больше 0")
		}

		value := uint64(in.Discount.Value)
		product.Discount = &entities.Discount{
			Type:     &discountType,
			Value:    &value,
			StartsAt: stringPtrOrNil(strings.TrimSpace(in.Discount.StartsAt)),
			EndsAt:   stringPtrOrNil(strings.TrimSpace(in.Discount.EndsAt)),
		}
	}

	return product, nil
}

func ValidateProductMedia(product *entities.Product) error {
	if product == nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "пустой продукт")
	}

	if product.MainImageURL == nil || strings.TrimSpace(*product.MainImageURL) == "" {
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "не передан main_image_url")
	}

	if len(product.Images) == 0 {
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "не передан список images")
	}

	for _, image := range product.Images {
		if strings.TrimSpace(image) == "" {
			return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "images содержит пустую ссылку")
		}
	}

	return nil
}

func ValidateProductUpdateImages(mainImageHeader *multipart.FileHeader, imageHeaders []*multipart.FileHeader) error {
	if mainImageHeader != nil {
		if _, err := ValidateImageExtension(mainImageHeader.Filename, []string{".jpg", ".jpeg", ".png", ".webp"}); err != nil {
			return err
		}
	}

	for _, header := range imageHeaders {
		if header == nil {
			return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "передан пустой заголовок картинки продукта")
		}
		if _, err := ValidateImageExtension(header.Filename, []string{".jpg", ".jpeg", ".png", ".webp"}); err != nil {
			return err
		}
	}

	return nil
}

func ValidateProductImages(mainImageHeader *multipart.FileHeader, imageHeaders []*multipart.FileHeader) error {
	switch {
	case mainImageHeader == nil:
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "не передана главная картинка продукта")
	case len(imageHeaders) == 0:
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "не переданы дополнительные картинки продукта")
	}

	if _, err := ValidateImageExtension(mainImageHeader.Filename, []string{".jpg", ".jpeg", ".png", ".webp"}); err != nil {
		return err
	}

	for _, header := range imageHeaders {
		if header == nil {
			return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "передан пустой заголовок картинки продукта")
		}
		if _, err := ValidateImageExtension(header.Filename, []string{".jpg", ".jpeg", ".png", ".webp"}); err != nil {
			return err
		}
	}

	return nil
}

func stringPtrOrNil(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func uint64Ptr(v uint64) *uint64 {
	return &v
}
