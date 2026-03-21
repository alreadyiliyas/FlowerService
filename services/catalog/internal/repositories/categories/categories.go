package repositories

import (
	"context"
	"errors"
	"time"

	convert "github.com/ilyas/flower/services/auth/internal/utils"
	"github.com/ilyas/flower/services/catalog/internal/dto"
	tnt "github.com/tarantool/go-tarantool/v2"
)

type tarantoolRepository struct {
	conn *tnt.Connection
}

func NewTarantoolRepository(conn *tnt.Connection) CategoryRepository {
	return &tarantoolRepository{conn: conn}
}

func (r *tarantoolRepository) List(ctx context.Context) ([]dto.Category, error) {
	return nil, nil
}

func (r *tarantoolRepository) Get(ctx context.Context, id uint64) (*dto.Category, error) {
	return nil, nil
}

func (r *tarantoolRepository) Create(ctx context.Context, in dto.Category) (*dto.Category, error) {
	req := tnt.NewCallRequest("create_category").
		Args([]interface{}{
			in.Name,
			in.Slug,
			in.Description,
			in.ImageURL,
		}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from create_category")
	}

	row, ok := resp[0].([]interface{})
	if !ok || len(row) < 7 {
		return nil, errors.New("invalid response payload")
	}
	id, err := convert.ToUint64(row[0], "id")
	if err != nil {
		return nil, err
	}
	name, err := convert.ToString(row[1], "name")
	if err != nil {
		return nil, err
	}
	slug, err := convert.ToString(row[2], "slug")
	if err != nil {
		return nil, err
	}
	description, err := convert.ToString(row[3], "description")
	if err != nil {
		return nil, err
	}
	image_url, err := convert.ToUint64(row[4], "image_url")
	if err != nil {
		return nil, err
	}
	createdUnix, err := convert.ToUint64(row[5], "created_at")
	if err != nil {
		return nil, err
	}
	updatedUnix, err := convert.ToUint64(row[6], "updated_at")
	if err != nil {
		return nil, err
	}

	createdAt := time.Unix(int64(createdUnix), 0)
	updatedAt := time.Unix(int64(updatedUnix), 0)

	return nil
}

func (r *tarantoolRepository) Update(ctx context.Context, id uint64, in dto.Category) (*dto.Category, error) {
	return nil, nil
}

func (r *tarantoolRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}
