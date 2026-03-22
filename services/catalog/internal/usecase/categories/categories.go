package usecase

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/dto"
	"github.com/ilyas/flower/services/catalog/internal/entities"
	repo "github.com/ilyas/flower/services/catalog/internal/repositories/categories"
	"github.com/ilyas/flower/services/catalog/internal/utils"
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

	items, err := uc.categories.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.Category, 0, len(items))
	for _, item := range items {
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
		result = append(result, category)
	}

	return result, nil
}

func (uc *categoriesUsecase) GetCategory(ctx context.Context, id uint64) (*dto.Category, error) {
	if uc.categories == nil {
		return nil, apperrors.ErrDB
	}

	item, err := uc.categories.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &dto.Category{}
	if item.ID != nil {
		resp.ID = *item.ID
	}
	if item.Name != nil {
		resp.Name = *item.Name
	}
	if item.Slug != nil {
		resp.Slug = *item.Slug
	}
	if item.Description != nil {
		resp.Description = *item.Description
	}
	if item.ImageURL != nil {
		resp.ImageURL = *item.ImageURL
	}
	if item.CreatedAt != nil {
		resp.CreatedAt = item.CreatedAt.Format(time.RFC3339)
	}
	if item.UpdatedAt != nil {
		resp.UpdatedAt = item.UpdatedAt.Format(time.RFC3339)
	}

	return resp, nil
}

func (uc *categoriesUsecase) CreateCategory(ctx context.Context, in dto.CreateCategoryRequest) (*dto.Category, error) {
	validCategory, err := utils.ValidateCategory(in)
	if err != nil {
		log.Printf("| usecase | create category error: %v", err)
		return nil, err
	}

	uploaded, err := utils.UploadImage(utils.UploadImageParams{
		File:         in.Image,
		Header:       in.ImageHeader,
		Dir:          "public/category",
		PublicPrefix: "/public/category",
		AllowedExt:   []string{".jpg", ".jpeg", ".png", ".webp"},
		FileNameSize: 16,
	})
	if err != nil {
		log.Printf("| usecase | create category upload failed: %v", err)
		return nil, apperrors.ErrDB
	}

	entity := entities.Category{
		Name:     validCategory.Name,
		Slug:     validCategory.Slug,
		ImageURL: &uploaded.PublicURL,
	}
	if *validCategory.Description != "" {
		entity.Description = validCategory.Description
	}

	created, err := uc.categories.Create(ctx, entity)
	if err != nil {
		utils.DeleteFileIfExists(uploaded.FullPath)
		log.Printf("create category repository failed: slug=%s err=%v", *validCategory.Slug, err)
		return nil, err
	}

	resp := &dto.Category{}
	if created.ID != nil {
		resp.ID = *created.ID
	}
	if created.Name != nil {
		resp.Name = *created.Name
	}
	if created.Slug != nil {
		resp.Slug = *created.Slug
	}
	if created.Description != nil {
		resp.Description = *created.Description
	}
	if created.ImageURL != nil {
		resp.ImageURL = *created.ImageURL
	}
	if created.CreatedAt != nil {
		resp.CreatedAt = created.CreatedAt.Format(time.RFC3339)
	}
	if created.UpdatedAt != nil {
		resp.UpdatedAt = created.UpdatedAt.Format(time.RFC3339)
	}

	log.Printf("create category success: id=%d slug=%s", resp.ID, resp.Slug)
	return resp, nil
}

func (uc *categoriesUsecase) UpdateCategory(ctx context.Context, id uint64, in dto.Category) (*dto.Category, error) {
	if uc.categories == nil {
		return nil, apperrors.ErrDB
	}

	entity := entities.Category{}
	if in.Name != "" {
		entity.Name = &in.Name
	}
	if in.Slug != "" {
		slug := strings.ToLower(strings.TrimSpace(in.Slug))
		entity.Slug = &slug
	}
	if in.Description != "" {
		entity.Description = &in.Description
	}
	if in.ImageURL != "" {
		entity.ImageURL = &in.ImageURL
	}

	item, err := uc.categories.Update(ctx, id, entity)
	if err != nil {
		return nil, err
	}

	resp := &dto.Category{}
	if item.ID != nil {
		resp.ID = *item.ID
	}
	if item.Name != nil {
		resp.Name = *item.Name
	}
	if item.Slug != nil {
		resp.Slug = *item.Slug
	}
	if item.Description != nil {
		resp.Description = *item.Description
	}
	if item.ImageURL != nil {
		resp.ImageURL = *item.ImageURL
	}
	if item.CreatedAt != nil {
		resp.CreatedAt = item.CreatedAt.Format(time.RFC3339)
	}
	if item.UpdatedAt != nil {
		resp.UpdatedAt = item.UpdatedAt.Format(time.RFC3339)
	}

	return resp, nil
}

func (uc *categoriesUsecase) DeleteCategory(ctx context.Context, id uint64) error {
	if uc.categories == nil {
		return apperrors.ErrDB
	}
	return uc.categories.Delete(ctx, id)
}
