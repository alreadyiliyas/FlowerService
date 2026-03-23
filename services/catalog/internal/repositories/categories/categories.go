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

func NewTarantoolRepository(conn *tnt.Connection) CategoryRepository {
	return &tarantoolRepository{conn: conn}
}

func (r *tarantoolRepository) List(ctx context.Context) ([]entities.Category, error) {
	req := tnt.NewCallRequest("list_categories").Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return []entities.Category{}, nil
	}

	return convert.ParseCategoryEntities(resp[0])
}

func (r *tarantoolRepository) Get(ctx context.Context, id uint64) (*entities.Category, error) {
	req := tnt.NewCallRequest("get_category").Args([]interface{}{id}).Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from get_category")
	}

	return convert.ParseCategoryEntity(resp[0])
}

func (r *tarantoolRepository) Create(ctx context.Context, in entities.Category) (*entities.Category, error) {
	req := tnt.NewCallRequest("create_category").
		Args([]interface{}{
			convert.ValueOrNull(in.Name),
			convert.ValueOrNull(in.Slug),
			convert.ValueOrNull(in.Description),
			convert.ValueOrNull(in.ImageURL),
		}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from create_category")
	}

	return convert.ParseCategoryEntity(resp[0])
}

func (r *tarantoolRepository) Update(ctx context.Context, id uint64, in entities.Category) (*entities.Category, error) {
	req := tnt.NewCallRequest("update_category").
		Args([]interface{}{
			id,
			convert.ValueOrNull(in.Name),
			convert.ValueOrNull(in.Slug),
			convert.ValueOrNull(in.Description),
			convert.ValueOrNull(in.ImageURL),
		}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from update_category")
	}

	return convert.ParseCategoryEntity(resp[0])
}

func (r *tarantoolRepository) Delete(ctx context.Context, id uint64) error {
	req := tnt.NewCallRequest("delete_category").Args([]interface{}{id}).Context(ctx)

	_, err := r.conn.Do(req).Get()
	if err != nil {
		return convert.MapTarantoolError(err)
	}
	return nil
}
