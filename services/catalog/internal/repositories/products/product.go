package repositories

import (
	"context"
	"errors"

	"github.com/ilyas/flower/services/catalog/internal/dto"
	"github.com/ilyas/flower/services/catalog/internal/entities"
	convert "github.com/ilyas/flower/services/catalog/internal/utils"
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

func (r *tarantoolRepository) Create(ctx context.Context, in entities.Product) (*entities.Product, error) {
	req := tnt.NewCallRequest("create_product").
		Args([]interface{}{
			*in.Name,
			convert.ValueOrNull(in.Description),
			*in.CategoryID,
			*in.SellerID,
			in.IsAvailable,
			*in.Currency,
			*in.MainImageURL,
			convert.BuildImagesArgs(in.Images),
			convert.BuildSizesArgs(in.Sizes),
			*in.PricePerStem,
			*in.MinStems,
			*in.MaxStems,
			convert.BuildCompositionArgs(in.Composition),
			convert.BuildDiscountArg(in.Discount),
		}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from create_product")
	}

	return convert.ParseProductEntity(resp[0])
}

func (r *tarantoolRepository) Update(ctx context.Context, id uint64, in entities.Product) (*entities.Product, error) {
	return nil, nil
}

func (r *tarantoolRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}
