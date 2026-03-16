package repositories

import (
	"context"

	"github.com/ilyas/flower/services/catalog/internal/dto"
)

type ProductRepository interface {
	List(ctx context.Context, filter dto.ProductFilter) (dto.PaginatedProducts, error)
	Get(ctx context.Context, id uint64) (*dto.Product, error)
	Create(ctx context.Context, in dto.Product) (*dto.Product, error)
	Update(ctx context.Context, id uint64, in dto.Product) (*dto.Product, error)
	Delete(ctx context.Context, id uint64) error
}
