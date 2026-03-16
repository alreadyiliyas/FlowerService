package usecase

import (
	"context"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/dto"

	repo "github.com/ilyas/flower/services/catalog/internal/repositories/categories"
)

type categoriesUsecase struct {
	categories repo.CategoryRepository
}

func NewCategoriesUsecase(categories repo.CategoryRepository) UsecaseCategories {
	return &categoriesUsecase{
		categories: categories,
	}
}

func (uc *categoriesUsecase) ListCategories(ctx context.Context) ([]dto.Category, error) {
	if uc.categories == nil {
		return nil, apperrors.ErrDB
	}
	return uc.categories.List(ctx)
}

func (uc *categoriesUsecase) GetCategory(ctx context.Context, id uint64) (*dto.Category, error) {
	if uc.categories == nil {
		return nil, apperrors.ErrDB
	}
	return uc.categories.Get(ctx, id)
}

func (uc *categoriesUsecase) CreateCategory(ctx context.Context, in dto.Category) (*dto.Category, error) {
	if uc.categories == nil {
		return nil, apperrors.ErrDB
	}
	return uc.categories.Create(ctx, in)
}

func (uc *categoriesUsecase) UpdateCategory(ctx context.Context, id uint64, in dto.Category) (*dto.Category, error) {
	if uc.categories == nil {
		return nil, apperrors.ErrDB
	}
	return uc.categories.Update(ctx, id, in)
}

func (uc *categoriesUsecase) DeleteCategory(ctx context.Context, id uint64) error {
	if uc.categories == nil {
		return apperrors.ErrDB
	}
	return uc.categories.Delete(ctx, id)
}
