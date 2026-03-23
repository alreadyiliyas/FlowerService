package repositories

import (
	"context"

	"github.com/ilyas/flower/services/catalog/internal/entities"
)

type ProductRepository interface {
	List(ctx context.Context, filter entities.ProductFilter) (entities.PaginatedProducts, error)
	Get(ctx context.Context, id uint64) (*entities.Product, error)
	Create(ctx context.Context, in entities.Product) (*entities.Product, error)
	Update(ctx context.Context, id uint64, in entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id uint64) error
}
