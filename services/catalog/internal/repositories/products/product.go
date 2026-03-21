package repositories

import (
	"context"

	"github.com/ilyas/flower/services/catalog/internal/dto"
	tnt "github.com/tarantool/go-tarantool/v2"
)

type tarantoolRepository struct {
	conn *tnt.Connection
}

func NewTarantoolRepository(conn *tnt.Connection) ProductRepository {
	return &tarantoolRepository{conn: conn}
}

func (r *tarantoolRepository) List(ctx context.Context, filter dto.ProductFilter) (dto.PaginatedProducts, error) {
	return dto.PaginatedProducts{}, nil
}

func (r *tarantoolRepository) Get(ctx context.Context, id uint64) (*dto.Product, error) {
	return nil, nil
}

func (r *tarantoolRepository) Create(ctx context.Context, in dto.Product) (*dto.Product, error) {
	return nil, nil
}

func (r *tarantoolRepository) Update(ctx context.Context, id uint64, in dto.Product) (*dto.Product, error) {
	return nil, nil
}

func (r *tarantoolRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}
