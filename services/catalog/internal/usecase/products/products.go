package usecase

import (
	"context"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/dto"
	repo "github.com/ilyas/flower/services/catalog/internal/repositories/products"
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

func (uc *productsUsecase) CreateProduct(ctx context.Context, in dto.Product) (*dto.Product, error) {
	if uc.products == nil {
		return nil, apperrors.ErrDB
	}
	return uc.products.Create(ctx, in)
}

func (uc *productsUsecase) UpdateProduct(ctx context.Context, id uint64, in dto.Product) (*dto.Product, error) {
	if uc.products == nil {
		return nil, apperrors.ErrDB
	}
	return uc.products.Update(ctx, id, in)
}

func (uc *productsUsecase) DeleteProduct(ctx context.Context, id uint64) error {
	if uc.products == nil {
		return apperrors.ErrDB
	}
	return uc.products.Delete(ctx, id)
}
