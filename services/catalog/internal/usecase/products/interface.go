package usecase

import (
	"context"

	"github.com/ilyas/flower/services/catalog/internal/dto"
)

type ProductUsecase interface {
	ListProducts(ctx context.Context, filter dto.ProductFilter) (dto.PaginatedProducts, error)
	GetProduct(ctx context.Context, id uint64) (*dto.Product, error)
	CreateProduct(ctx context.Context, in dto.CreateProductRequest) (*dto.Product, error)
	UpdateProduct(ctx context.Context, id uint64, in dto.UpdateProductRequest) (*dto.Product, error)
	DeleteProduct(ctx context.Context, id uint64, in dto.DeleteProductRequest) error
}
