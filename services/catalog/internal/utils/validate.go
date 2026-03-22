package utils

import (
	"fmt"
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
