package usecase

import (
	"context"
	"log"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/dto"
	repo "github.com/ilyas/flower/services/catalog/internal/repositories/products"
	"github.com/ilyas/flower/services/catalog/internal/utils"
)

type productsUsecase struct {
	products repo.ProductRepository
}

func NewproductsUsecase(products repo.ProductRepository) ProductUsecase {
	return &productsUsecase{
		products: products,
	}
}

func (uc *productsUsecase) ListProducts(ctx context.Context, filter dto.ProductFilter) (dto.PaginatedProducts, error) {
	if uc.products == nil {
		return dto.PaginatedProducts{}, apperrors.ErrDB
	}
	return uc.products.List(ctx, filter)
}

func (uc *productsUsecase) GetProduct(ctx context.Context, id uint64) (*dto.Product, error) {
	if uc.products == nil {
		return nil, apperrors.ErrDB
	}
	return uc.products.Get(ctx, id)
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

	return utils.MapProductEntityToDTO(*resp), nil
}

func (uc *productsUsecase) UpdateProduct(ctx context.Context, id uint64, in dto.Product) (*dto.Product, error) {
	if uc.products == nil {
		return nil, apperrors.ErrDB
	}
	return nil, nil
}

func (uc *productsUsecase) DeleteProduct(ctx context.Context, id uint64) error {
	if uc.products == nil {
		return apperrors.ErrDB
	}
	return uc.products.Delete(ctx, id)
}
