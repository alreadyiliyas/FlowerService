package repositories

import (
	"context"
	"errors"

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

func (r *tarantoolRepository) List(ctx context.Context, filter entities.ProductFilter) (entities.PaginatedProducts, error) {
	req := tnt.NewCallRequest("list_products").
		Args([]interface{}{
			map[string]interface{}{
				"category_id":  convert.Uint64OrNull(filter.CategoryID),
				"seller_id":    convert.Uint64OrNull(filter.SellerID),
				"price_min":    convert.Uint64OrNull(filter.PriceMin),
				"price_max":    convert.Uint64OrNull(filter.PriceMax),
				"size":         convert.ValueOrNull(filter.Size),
				"is_available": convert.BoolOrNull(filter.IsAvailable),
				"page":         filter.Page,
				"page_size":    filter.PageSize,
			},
		}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return entities.PaginatedProducts{}, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return entities.PaginatedProducts{}, errors.New("empty response from list_products")
	}

	return convert.ParsePaginatedProductEntities(resp[0])
}

func (r *tarantoolRepository) Get(ctx context.Context, id uint64) (*entities.Product, error) {
	req := tnt.NewCallRequest("get_product").Args([]interface{}{id}).Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from get_product")
	}

	return convert.ParseProductEntity(resp[0])
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
	req := tnt.NewCallRequest("update_product").
		Args([]interface{}{
			id,
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
			in.Version,
		}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from update_product")
	}

	return convert.ParseProductEntity(resp[0])
}

func (r *tarantoolRepository) Delete(ctx context.Context, id uint64) error {
	req := tnt.NewCallRequest("delete_product").Args([]interface{}{id}).Context(ctx)

	_, err := r.conn.Do(req).Get()
	if err != nil {
		return convert.MapTarantoolError(err)
	}
	return nil
}
