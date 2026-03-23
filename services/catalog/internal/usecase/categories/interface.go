package usecase

import (
	"context"

	"github.com/ilyas/flower/services/catalog/internal/dto"
)

type UsecaseCategories interface {
	ListCategories(ctx context.Context) ([]dto.Category, error)
	GetCategory(ctx context.Context, id uint64) (*dto.Category, error)
	CreateCategory(ctx context.Context, in dto.CreateCategoryRequest) (*dto.Category, error)
	UpdateCategory(ctx context.Context, id uint64, in dto.UpdateCategoryRequest) (*dto.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
}
