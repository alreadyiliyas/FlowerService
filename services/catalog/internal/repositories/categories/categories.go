package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

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

	rows, ok := resp[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response payload type: %T", resp[0])
	}

	result := make([]entities.Category, 0, len(rows))
	for _, raw := range rows {
		row, ok := raw.(map[interface{}]interface{})
		if !ok {
			if mapped, ok := raw.(map[string]interface{}); ok {
				row = make(map[interface{}]interface{}, len(mapped))
				for k, v := range mapped {
					row[k] = v
				}
			} else {
				return nil, fmt.Errorf("invalid category row type: %T", raw)
			}
		}

		id, err := convert.ToUint64(row["id"], "id")
		if err != nil {
			return nil, err
		}
		name, err := convert.ToString(row["name"], "name")
		if err != nil {
			return nil, err
		}
		slug, err := convert.ToString(row["slug"], "slug")
		if err != nil {
			return nil, err
		}
		description, err := convert.ToStringNullable(row["description"], "description")
		if err != nil {
			return nil, err
		}
		imageURL, err := convert.ToStringNullable(row["image_url"], "image_url")
		if err != nil {
			return nil, err
		}
		createdUnix, err := convert.ToUint64(row["created_at"], "created_at")
		if err != nil {
			return nil, err
		}
		updatedUnix, err := convert.ToUint64(row["updated_at"], "updated_at")
		if err != nil {
			return nil, err
		}

		createdAt := time.Unix(int64(createdUnix), 0)
		updatedAt := time.Unix(int64(updatedUnix), 0)

		result = append(result, entities.Category{
			ID:          &id,
			Name:        &name,
			Slug:        &slug,
			Description: description,
			ImageURL:    imageURL,
			CreatedAt:   &createdAt,
			UpdatedAt:   &updatedAt,
		})
	}

	return result, nil
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

	row, ok := resp[0].(map[interface{}]interface{})
	if !ok {
		if mapped, ok := resp[0].(map[string]interface{}); ok {
			row = make(map[interface{}]interface{}, len(mapped))
			for k, v := range mapped {
				row[k] = v
			}
		} else {
			return nil, fmt.Errorf("invalid response payload type: %T", resp[0])
		}
	}

	categoryID, err := convert.ToUint64(row["id"], "id")
	if err != nil {
		return nil, err
	}
	name, err := convert.ToString(row["name"], "name")
	if err != nil {
		return nil, err
	}
	slug, err := convert.ToString(row["slug"], "slug")
	if err != nil {
		return nil, err
	}
	description, err := convert.ToStringNullable(row["description"], "description")
	if err != nil {
		return nil, err
	}
	imageURL, err := convert.ToStringNullable(row["image_url"], "image_url")
	if err != nil {
		return nil, err
	}
	createdUnix, err := convert.ToUint64(row["created_at"], "created_at")
	if err != nil {
		return nil, err
	}
	updatedUnix, err := convert.ToUint64(row["updated_at"], "updated_at")
	if err != nil {
		return nil, err
	}

	createdAt := time.Unix(int64(createdUnix), 0)
	updatedAt := time.Unix(int64(updatedUnix), 0)

	return &entities.Category{
		ID:          &categoryID,
		Name:        &name,
		Slug:        &slug,
		Description: description,
		ImageURL:    imageURL,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}, nil
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

	row, ok := resp[0].(map[interface{}]interface{})
	if !ok {
		if mapped, ok := resp[0].(map[string]interface{}); ok {
			row = make(map[interface{}]interface{}, len(mapped))
			for k, v := range mapped {
				row[k] = v
			}
		} else {
			return nil, fmt.Errorf("invalid response payload type: %T", resp[0])
		}
	}

	id, err := convert.ToUint64(row["id"], "id")
	if err != nil {
		return nil, err
	}
	name, err := convert.ToString(row["name"], "name")
	if err != nil {
		return nil, err
	}
	slug, err := convert.ToString(row["slug"], "slug")
	if err != nil {
		return nil, err
	}
	description, err := convert.ToStringNullable(row["description"], "description")
	if err != nil {
		return nil, err
	}
	imageURL, err := convert.ToStringNullable(row["image_url"], "image_url")
	if err != nil {
		return nil, err
	}
	createdUnix, err := convert.ToUint64(row["created_at"], "created_at")
	if err != nil {
		return nil, err
	}
	updatedUnix, err := convert.ToUint64(row["updated_at"], "updated_at")
	if err != nil {
		return nil, err
	}
	createdAt := time.Unix(int64(createdUnix), 0)
	updatedAt := time.Unix(int64(updatedUnix), 0)

	return &entities.Category{
		ID:          &id,
		Name:        &name,
		Slug:        &slug,
		Description: description,
		ImageURL:    imageURL,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}, nil
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

	row, ok := resp[0].(map[interface{}]interface{})
	if !ok {
		if mapped, ok := resp[0].(map[string]interface{}); ok {
			row = make(map[interface{}]interface{}, len(mapped))
			for k, v := range mapped {
				row[k] = v
			}
		} else {
			return nil, fmt.Errorf("invalid response payload type: %T", resp[0])
		}
	}

	categoryID, err := convert.ToUint64(row["id"], "id")
	if err != nil {
		return nil, err
	}
	name, err := convert.ToString(row["name"], "name")
	if err != nil {
		return nil, err
	}
	slug, err := convert.ToString(row["slug"], "slug")
	if err != nil {
		return nil, err
	}
	description, err := convert.ToStringNullable(row["description"], "description")
	if err != nil {
		return nil, err
	}
	imageURL, err := convert.ToStringNullable(row["image_url"], "image_url")
	if err != nil {
		return nil, err
	}
	createdUnix, err := convert.ToUint64(row["created_at"], "created_at")
	if err != nil {
		return nil, err
	}
	updatedUnix, err := convert.ToUint64(row["updated_at"], "updated_at")
	if err != nil {
		return nil, err
	}
	createdAt := time.Unix(int64(createdUnix), 0)
	updatedAt := time.Unix(int64(updatedUnix), 0)

	return &entities.Category{
		ID:          &categoryID,
		Name:        &name,
		Slug:        &slug,
		Description: description,
		ImageURL:    imageURL,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}, nil
}

func (r *tarantoolRepository) Delete(ctx context.Context, id uint64) error {
	req := tnt.NewCallRequest("delete_category").Args([]interface{}{id}).Context(ctx)

	_, err := r.conn.Do(req).Get()
	if err != nil {
		return convert.MapTarantoolError(err)
	}
	return nil
}
