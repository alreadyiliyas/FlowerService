package usecase

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/dto"
	"github.com/ilyas/flower/services/catalog/internal/entities"
	cacherepo "github.com/ilyas/flower/services/catalog/internal/repositories/cache"
	repo "github.com/ilyas/flower/services/catalog/internal/repositories/products"
	"github.com/ilyas/flower/services/catalog/internal/utils"
	"github.com/redis/go-redis/v9"
)

// ToDo вынести env
const productCacheTTL = 15 * time.Minute

type productsUsecase struct {
	products repo.ProductRepository
	cache    cacherepo.CacheRepository
}

func NewproductsUsecase(products repo.ProductRepository, cache cacherepo.CacheRepository) ProductUsecase {
	return &productsUsecase{
		products: products,
		cache:    cache,
	}
}

func (uc *productsUsecase) ListProducts(ctx context.Context, filter dto.ProductFilter) (dto.PaginatedProducts, error) {
	if uc.products == nil {
		return dto.PaginatedProducts{}, apperrors.ErrDB
	}

	listVersion := uc.getProductsListVersion(ctx)
	listKey := utils.BuildProductsListKey(filter, listVersion)

	if uc.cache != nil {
		cached, err := uc.cache.Get(ctx, listKey)
		switch {
		case err == nil:
			var items dto.PaginatedProducts
			if err := utils.UnmarshalFromString(cached, &items); err == nil {
				return items, nil
			}
		case errors.Is(err, redis.Nil):
		default:
			log.Printf("| usecase | list products cache get failed: %v", err)
		}
	}

	page, err := uc.products.List(ctx, utils.MapProductFilterToEntity(filter))
	if err != nil {
		log.Printf("| usecase | list products repository failed: %v", err)
		return dto.PaginatedProducts{}, err
	}

	resp := utils.MapPaginatedProductsToDTO(page)
	if uc.cache != nil {
		if raw, err := utils.MarshalToString(resp); err == nil {
			_ = uc.cache.Set(ctx, listKey, raw, productCacheTTL)
		}
		for _, item := range resp.Items {
			if raw, err := utils.MarshalToString(item); err == nil {
				_ = uc.cache.Set(ctx, utils.BuildProductKey(item.ID), raw, productCacheTTL)
			}
		}
	}

	return resp, nil
}

func (uc *productsUsecase) GetProduct(ctx context.Context, id uint64) (*dto.Product, error) {
	if uc.products == nil {
		return nil, apperrors.ErrDB
	}

	itemKey := utils.BuildProductKey(id)
	if uc.cache != nil {
		cached, err := uc.cache.Get(ctx, itemKey)
		switch {
		case err == nil:
			var item dto.Product
			if err := utils.UnmarshalFromString(cached, &item); err == nil {
				return &item, nil
			}
		case errors.Is(err, redis.Nil):
		default:
			log.Printf("| usecase | get product cache get failed: %v", err)
		}
	}

	item, err := uc.products.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := utils.MapProductEntityToDTO(*item)
	if uc.cache != nil {
		if raw, err := utils.MarshalToString(resp); err == nil {
			_ = uc.cache.Set(ctx, itemKey, raw, productCacheTTL)
		}
	}

	return resp, nil
}

func (uc *productsUsecase) CreateProduct(ctx context.Context, in dto.CreateProductRequest) (*dto.Product, error) {
	if uc.products == nil {
		return nil, apperrors.ErrDB
	}

	product, err := utils.ValidateProduct(in.Product)
	if err != nil {
		log.Printf("| usecase | create product payload validation failed: %v", err)
		return nil, err
	}

	if err := utils.ValidateProductImages(in.MainImageHeader, in.ImageHeaders); err != nil {
		log.Printf("| usecase | create product images validation failed: %v", err)
		return nil, err
	}

	if err := uc.validateProductTypeRole(in.TypeUserID, in.TypeRole); err != nil {
		return nil, err
	}

	product.SellerID = &in.TypeUserID

	defer func() {
		if in.MainImage != nil {
			_ = in.MainImage.Close()
		}
	}()

	mainImage, err := utils.UploadImage(utils.UploadImageParams{
		File:         in.MainImage,
		Header:       in.MainImageHeader,
		Dir:          "public/products",
		PublicPrefix: "/public/products",
		AllowedExt:   []string{".jpg", ".jpeg", ".png", ".webp"},
		FileNameSize: 16,
	})
	if err != nil {
		log.Printf("| usecase | create product main image upload failed: %v", err)
		return nil, err
	}

	uploadedImages, err := utils.UploadImages(utils.UploadImagesParams{
		Files:        in.Images,
		Headers:      in.ImageHeaders,
		Dir:          "public/products",
		PublicPrefix: "/public/products",
		AllowedExt:   []string{".jpg", ".jpeg", ".png", ".webp"},
		FileNameSize: 16,
	})
	if err != nil {
		utils.DeleteFileIfExists(mainImage.FullPath)
		log.Printf("| usecase | create product images upload failed: %v", err)
		return nil, err
	}

	product.MainImageURL = &mainImage.PublicURL
	product.Images = make([]string, 0, len(uploadedImages))
	for _, item := range uploadedImages {
		product.Images = append(product.Images, item.PublicURL)
	}

	if err := utils.ValidateProductMedia(product); err != nil {
		utils.DeleteFileIfExists(mainImage.FullPath)
		utils.DeleteUploadedFiles(uploadedImages)
		log.Printf("| usecase | create product final media validation failed: %v", err)
		return nil, err
	}

	resp, err := uc.products.Create(ctx, *product)
	if err != nil {
		log.Printf("| usecase | create product repositories failed: %v", err)
		utils.DeleteFileIfExists(mainImage.FullPath)
		utils.DeleteUploadedFiles(uploadedImages)
		return nil, err
	}

	dtoResp := utils.MapProductEntityToDTO(*resp)
	uc.refreshProductCache(ctx, dtoResp)
	return dtoResp, nil
}

func (uc *productsUsecase) UpdateProduct(ctx context.Context, id uint64, in dto.UpdateProductRequest) (*dto.Product, error) {
	if uc.products == nil {
		return nil, apperrors.ErrDB
	}

	current, err := uc.products.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := uc.ensureProductAccess(current, in.TypeUserID, in.TypeRole); err != nil {
		return nil, err
	}

	product, err := utils.ValidateProduct(in.Product)
	if err != nil {
		log.Printf("| usecase | update product payload validation failed: %v", err)
		return nil, err
	}

	if err := utils.ValidateProductUpdateImages(in.MainImageHeader, in.ImageHeaders); err != nil {
		log.Printf("| usecase | update product image validation failed: %v", err)
		return nil, err
	}

	if in.MainImage != nil {
		defer func() { _ = in.MainImage.Close() }()
	}

	product.ID = current.ID
	product.Version = current.Version
	product.SellerID = current.SellerID

	oldMainImageURL := ""
	if current.MainImageURL != nil {
		oldMainImageURL = *current.MainImageURL
	}
	oldImages := append([]string(nil), current.Images...)

	var newMainImage *utils.UploadedFile
	if in.MainImage != nil && in.MainImageHeader != nil {
		newMainImage, err = utils.UploadImage(utils.UploadImageParams{
			File:         in.MainImage,
			Header:       in.MainImageHeader,
			Dir:          "public/products",
			PublicPrefix: "/public/products",
			AllowedExt:   []string{".jpg", ".jpeg", ".png", ".webp"},
			FileNameSize: 16,
		})
		if err != nil {
			log.Printf("| usecase | update product main image upload failed: %v", err)
			return nil, err
		}
		product.MainImageURL = &newMainImage.PublicURL
	} else {
		product.MainImageURL = current.MainImageURL
	}

	var uploadedImages []utils.UploadedFile
	if len(in.ImageHeaders) > 0 {
		uploadedImages, err = utils.UploadImages(utils.UploadImagesParams{
			Files:        in.Images,
			Headers:      in.ImageHeaders,
			Dir:          "public/products",
			PublicPrefix: "/public/products",
			AllowedExt:   []string{".jpg", ".jpeg", ".png", ".webp"},
			FileNameSize: 16,
		})
		if err != nil {
			if newMainImage != nil {
				utils.DeleteFileIfExists(newMainImage.FullPath)
			}
			log.Printf("| usecase | update product images upload failed: %v", err)
			return nil, err
		}

		product.Images = make([]string, 0, len(uploadedImages))
		for _, item := range uploadedImages {
			product.Images = append(product.Images, item.PublicURL)
		}
	} else {
		product.Images = append([]string(nil), current.Images...)
	}

	if err := utils.ValidateProductMedia(product); err != nil {
		if newMainImage != nil {
			utils.DeleteFileIfExists(newMainImage.FullPath)
		}
		utils.DeleteUploadedFiles(uploadedImages)
		log.Printf("| usecase | update product final media validation failed: %v", err)
		return nil, err
	}

	updated, err := uc.products.Update(ctx, id, *product)
	if err != nil {
		if newMainImage != nil {
			utils.DeleteFileIfExists(newMainImage.FullPath)
		}
		utils.DeleteUploadedFiles(uploadedImages)
		log.Printf("| usecase | update product repositories failed: %v", err)
		return nil, err
	}

	if newMainImage != nil && oldMainImageURL != "" {
		utils.DeletePublicFile(oldMainImageURL)
	}
	if len(uploadedImages) > 0 {
		utils.DeletePublicFiles(oldImages)
	}

	dtoResp := utils.MapProductEntityToDTO(*updated)
	uc.refreshProductCache(ctx, dtoResp)
	return dtoResp, nil
}

func (uc *productsUsecase) DeleteProduct(ctx context.Context, id uint64, in dto.DeleteProductRequest) error {
	if uc.products == nil {
		return apperrors.ErrDB
	}

	current, err := uc.products.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.ensureProductAccess(current, in.TypeUserID, in.TypeRole); err != nil {
		return err
	}

	if err := uc.products.Delete(ctx, id); err != nil {
		return err
	}

	if current.MainImageURL != nil {
		utils.DeletePublicFile(*current.MainImageURL)
	}
	utils.DeletePublicFiles(current.Images)

	if uc.cache != nil {
		_ = uc.cache.Del(ctx, utils.BuildProductKey(id))
		uc.invalidateProductsListCache(ctx)
	}

	return nil
}

func (uc *productsUsecase) getProductsListVersion(ctx context.Context) string {
	if uc.cache == nil {
		return "1"
	}

	version, err := uc.cache.Get(ctx, utils.BuildProductsListVersionKey())
	switch {
	case err == nil && version != "":
		return version
	case errors.Is(err, redis.Nil):
		return "1"
	default:
		if err != nil {
			log.Printf("| usecase | get products list version failed: %v", err)
		}
		return "1"
	}
}

func (uc *productsUsecase) invalidateProductsListCache(ctx context.Context) {
	if uc.cache == nil {
		return
	}

	version := strconv.FormatInt(time.Now().UnixNano(), 10)
	_ = uc.cache.Set(ctx, utils.BuildProductsListVersionKey(), version, productCacheTTL)
}

func (uc *productsUsecase) refreshProductCache(ctx context.Context, item *dto.Product) {
	if uc.cache == nil || item == nil {
		return
	}

	uc.invalidateProductsListCache(ctx)
	if raw, err := utils.MarshalToString(item); err == nil {
		_ = uc.cache.Set(ctx, utils.BuildProductKey(item.ID), raw, productCacheTTL)
	}
}

func (uc *productsUsecase) validateProductTypeRole(actorUserID uint64, actorRole string) error {
	if actorUserID == 0 {
		return apperrors.ErrUnauthorized
	}

	switch strings.TrimSpace(actorRole) {
	case "seller", "moderator":
		return nil
	default:
		return apperrors.ErrForbidden
	}
}

func (uc *productsUsecase) ensureProductAccess(product *entities.Product, UserID uint64, typeRole string) error {
	if err := uc.validateProductTypeRole(UserID, typeRole); err != nil {
		return err
	}

	if strings.TrimSpace(typeRole) == "moderator" {
		return nil
	}

	if product == nil || product.SellerID == nil {
		return apperrors.ErrForbidden
	}

	if *product.SellerID != UserID {
		return apperrors.ErrForbidden
	}

	return nil
}
