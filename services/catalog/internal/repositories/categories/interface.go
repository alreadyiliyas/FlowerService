package repositories

import (
	"context"

	"github.com/ilyas/flower/services/catalog/internal/dto"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]dto.Category, error)
	Get(ctx context.Context, id uint64) (*dto.Category, error)
	Create(ctx context.Context, in dto.Category) (*dto.Category, error)
	Update(ctx context.Context, id uint64, in dto.Category) (*dto.Category, error)
	Delete(ctx context.Context, id uint64) error
}
