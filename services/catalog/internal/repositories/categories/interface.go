package repositories

import (
	"context"

	"github.com/ilyas/flower/services/catalog/internal/entities"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]entities.Category, error)
	Get(ctx context.Context, id uint64) (*entities.Category, error)
	Create(ctx context.Context, in entities.Category) (*entities.Category, error)
	Update(ctx context.Context, id uint64, in entities.Category) (*entities.Category, error)
	Delete(ctx context.Context, id uint64) error
}
