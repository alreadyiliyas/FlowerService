package usecase

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/dto"
	"github.com/ilyas/flower/services/catalog/internal/entities"
	cacherepo "github.com/ilyas/flower/services/catalog/internal/repositories/cache"
	repo "github.com/ilyas/flower/services/catalog/internal/repositories/categories"
	"github.com/ilyas/flower/services/catalog/internal/utils"
	"github.com/redis/go-redis/v9"
)

// ToDo env
const categoryCacheTTL = 15 * time.Minute

type categoriesUsecase struct {
	categories repo.CategoryRepository
	cache      cacherepo.CacheRepository
}

func NewCategoriesUsecase(categories repo.CategoryRepository, cache cacherepo.CacheRepository) UsecaseCategories {
	return &categoriesUsecase{
		categories: categories,
		cache:      cache,
	}
}

func (uc *categoriesUsecase) ListCategories(ctx context.Context) ([]dto.Category, error) {
	if uc.categories == nil {
		return nil, apperrors.ErrDB
	}

	listKey := utils.BuildCategoriesListKey()
	if uc.cache != nil {
		cached, err := uc.cache.Get(ctx, listKey)
		switch {
		case err == nil:
			var items []dto.Category
			if err := utils.UnmarshalFromString(cached, &items); err == nil {
				return items, nil
			}
		case errors.Is(err, redis.Nil):
		default:
			log.Printf("| usecase | list categories cache get failed: %v", err)
		}
	}

	items, err := uc.categories.List(ctx)
	if err != nil {
		log.Printf("| usecase | list categories get failed: %v", err)
		return nil, err
	}

	result := make([]dto.Category, 0, len(items))
	for _, item := range items {
		category := utils.MapCategoryToDTO(item)
		result = append(result, category)
		if uc.cache != nil && item.ID != nil {
			if raw, err := utils.MarshalToString(category); err == nil {
				_ = uc.cache.Set(ctx, utils.BuildCategoryKey(*item.ID), raw, categoryCacheTTL)
			}
		}
	}

	if uc.cache != nil {
		if raw, err := utils.MarshalToString(result); err == nil {
			_ = uc.cache.Set(ctx, listKey, raw, categoryCacheTTL)
		}
	}

	return result, nil
}

func (uc *categoriesUsecase) GetCategory(ctx context.Context, id uint64) (*dto.Category, error) {
	if uc.categories == nil {
		return nil, apperrors.ErrDB
	}

	itemKey := utils.BuildCategoryKey(id)
	if uc.cache != nil {
		cached, err := uc.cache.Get(ctx, itemKey)
		switch {
		case err == nil:
			var item dto.Category
			if err := utils.UnmarshalFromString(cached, &item); err == nil {
				return &item, nil
			}
		case errors.Is(err, redis.Nil):
		default:
			log.Printf("| usecase | get category cache get failed: %v", err)
		}
	}

	item, err := uc.categories.Get(ctx, id)
	if err != nil {
		log.Printf("| usecase | failed to get category: %v", err)
		return nil, err
	}

	resp := utils.MapCategoryToDTO(*item)
	if uc.cache != nil {
		if raw, err := utils.MarshalToString(resp); err == nil {
			_ = uc.cache.Set(ctx, itemKey, raw, categoryCacheTTL)
		}
	}

	return &resp, nil
}

func (uc *categoriesUsecase) CreateCategory(ctx context.Context, in dto.CreateCategoryRequest) (*dto.Category, error) {
	validCategory, err := utils.ValidateCategory(in)
	if err != nil {
		log.Printf("| usecase | create category error: %v", err)
		return nil, err
	}
	defer in.Image.Close()

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
		if errors.Is(err, apperrors.ErrInvalidInput) {
			return nil, err
		}
		return nil, apperrors.ErrDB
	}

	entity := entities.Category{
		Name:     validCategory.Name,
		Slug:     validCategory.Slug,
		ImageURL: &uploaded.PublicURL,
	}
	if validCategory.Description != nil && *validCategory.Description != "" {
		entity.Description = validCategory.Description
	}

	created, err := uc.categories.Create(ctx, entity)
	if err != nil {
		utils.DeleteFileIfExists(uploaded.FullPath)
		log.Printf("| usecase | create category repository failed: slug=%s err=%v", *validCategory.Slug, err)
		return nil, err
	}

	resp := utils.MapCategoryToDTO(*created)
	uc.refreshCategoryCache(ctx, &resp)

	return &resp, nil
}

func (uc *categoriesUsecase) UpdateCategory(ctx context.Context, id uint64, in dto.UpdateCategoryRequest) (*dto.Category, error) {
	if uc.categories == nil {
		log.Printf("| usecase | update category categoriesUsecase is nil")
		return nil, apperrors.ErrDB
	}

	current, err := uc.GetCategory(ctx, id)
	if err != nil {
		log.Printf("| usecase | failed to get category: %v", err)
		return nil, err
	}

	validCategory, err := utils.ValidateCategoryUpdate(in)
	if err != nil {
		log.Printf("| usecase | update category validation failed: %v", err)
		return nil, err
	}

	name := current.Name
	if validCategory.Name != nil {
		name = *validCategory.Name
	}
	slug := current.Slug
	if validCategory.Slug != nil {
		slug = *validCategory.Slug
	}
	description := current.Description
	if validCategory.Description != nil {
		description = *validCategory.Description
	}
	imageURL := current.ImageURL

	if in.Image != nil && in.ImageHeader != nil {
		defer in.Image.Close()
		if imageURL != "" {
			replaced, err := utils.ReplaceImage(utils.ReplaceImageParams{
				File:              in.Image,
				Header:            in.ImageHeader,
				ExistingPublicURL: imageURL,
				AllowedExt:        []string{".jpg", ".jpeg", ".png", ".webp"},
			})
			if err != nil {
				if errors.Is(err, apperrors.ErrInvalidInput) {
					return nil, err
				}
				return nil, apperrors.ErrDB
			}
			imageURL = replaced.PublicURL
		} else {
			uploaded, err := utils.UploadImage(utils.UploadImageParams{
				File:         in.Image,
				Header:       in.ImageHeader,
				Dir:          "public/category",
				PublicPrefix: "/public/category",
				AllowedExt:   []string{".jpg", ".jpeg", ".png", ".webp"},
				FileNameSize: 16,
			})
			if err != nil {
				if errors.Is(err, apperrors.ErrInvalidInput) {
					return nil, err
				}
				return nil, apperrors.ErrDB
			}
			imageURL = uploaded.PublicURL
		}
	}

	entity := entities.Category{
		Name:     &name,
		Slug:     &slug,
		ImageURL: &imageURL,
	}
	if description != "" {
		entity.Description = &description
	}

	updated, err := uc.categories.Update(ctx, id, entity)
	if err != nil {
		log.Printf("| usecase | failed to update category: %v", err)
		return nil, err
	}

	resp := utils.MapCategoryToDTO(*updated)
	uc.refreshCategoryCache(ctx, &resp)
	return &resp, nil
}

func (uc *categoriesUsecase) DeleteCategory(ctx context.Context, id uint64) error {
	if uc.categories == nil {
		log.Printf("| usecase | delete category categoriesUsecase is nil")
		return apperrors.ErrDB
	}

	current, err := uc.GetCategory(ctx, id)
	if err != nil {
		log.Printf("| usecase | failed to get category: %v", err)
		return err
	}

	if err := uc.categories.Delete(ctx, id); err != nil {
		log.Printf("| usecase | failed to delete category: %v", err)
		return err
	}

	if current.ImageURL != "" {
		utils.DeleteFileIfExists(strings.TrimPrefix(current.ImageURL, "/"))
	}
	if uc.cache != nil {
		_ = uc.cache.Del(ctx, utils.BuildCategoryKey(id))
		_ = uc.cache.Del(ctx, utils.BuildCategoriesListKey())
	}

	log.Printf("| usecase | delete category successful : %v", err)
	return nil
}

func (uc *categoriesUsecase) refreshCategoryCache(ctx context.Context, item *dto.Category) {
	if uc.cache == nil || item == nil {
		log.Printf("| usecase | failed cache or item is nil")
		return
	}

	_ = uc.cache.Del(ctx, utils.BuildCategoriesListKey())
	if raw, err := utils.MarshalToString(item); err == nil {
		_ = uc.cache.Set(ctx, utils.BuildCategoryKey(item.ID), raw, categoryCacheTTL)
	}
}
