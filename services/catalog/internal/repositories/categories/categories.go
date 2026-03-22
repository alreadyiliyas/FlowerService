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
	return nil, nil
}

func (r *tarantoolRepository) Get(ctx context.Context, id uint64) (*entities.Category, error) {
	return nil, nil
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

	var (
		id          uint64
		name        string
		slug        string
		description *string
		imageURL    *string
		createdAt   time.Time
		updatedAt   time.Time
	)

	switch row := resp[0].(type) {
	case []interface{}:
		if len(row) < 7 {
			return nil, errors.New("invalid response payload")
		}

		id, err = convert.ToUint64(row[0], "id")
		if err != nil {
			return nil, err
		}
		name, err = convert.ToString(row[1], "name")
		if err != nil {
			return nil, err
		}
		slug, err = convert.ToString(row[2], "slug")
		if err != nil {
			return nil, err
		}
		description, err = convert.ToStringNullable(row[3], "description")
		if err != nil {
			return nil, err
		}
		imageURL, err = convert.ToStringNullable(row[4], "image_url")
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
		createdAt = time.Unix(int64(createdUnix), 0)
		updatedAt = time.Unix(int64(updatedUnix), 0)
	case map[interface{}]interface{}:
		idRaw, ok := row["id"]
		if !ok {
			return nil, errors.New("missing id")
		}
		nameRaw, ok := row["name"]
		if !ok {
			return nil, errors.New("missing name")
		}
		slugRaw, ok := row["slug"]
		if !ok {
			return nil, errors.New("missing slug")
		}
		createdRaw, ok := row["created_at"]
		if !ok {
			return nil, errors.New("missing created_at")
		}
		updatedRaw, ok := row["updated_at"]
		if !ok {
			return nil, errors.New("missing updated_at")
		}

		id, err = convert.ToUint64(idRaw, "id")
		if err != nil {
			return nil, err
		}
		name, err = convert.ToString(nameRaw, "name")
		if err != nil {
			return nil, err
		}
		slug, err = convert.ToString(slugRaw, "slug")
		if err != nil {
			return nil, err
		}
		if raw, ok := row["description"]; ok && raw != nil {
			description, err = convert.ToStringNullable(raw, "description")
			if err != nil {
				return nil, err
			}
		}
		if raw, ok := row["image_url"]; ok && raw != nil {
			imageURL, err = convert.ToStringNullable(raw, "image_url")
			if err != nil {
				return nil, err
			}
		}
		createdUnix, err := convert.ToUint64(createdRaw, "created_at")
		if err != nil {
			return nil, err
		}
		updatedUnix, err := convert.ToUint64(updatedRaw, "updated_at")
		if err != nil {
			return nil, err
		}
		createdAt = time.Unix(int64(createdUnix), 0)
		updatedAt = time.Unix(int64(updatedUnix), 0)
	case map[string]interface{}:
		idRaw, ok := row["id"]
		if !ok {
			return nil, errors.New("missing id")
		}
		nameRaw, ok := row["name"]
		if !ok {
			return nil, errors.New("missing name")
		}
		slugRaw, ok := row["slug"]
		if !ok {
			return nil, errors.New("missing slug")
		}
		createdRaw, ok := row["created_at"]
		if !ok {
			return nil, errors.New("missing created_at")
		}
		updatedRaw, ok := row["updated_at"]
		if !ok {
			return nil, errors.New("missing updated_at")
		}

		id, err = convert.ToUint64(idRaw, "id")
		if err != nil {
			return nil, err
		}
		name, err = convert.ToString(nameRaw, "name")
		if err != nil {
			return nil, err
		}
		slug, err = convert.ToString(slugRaw, "slug")
		if err != nil {
			return nil, err
		}
		if raw, ok := row["description"]; ok && raw != nil {
			description, err = convert.ToStringNullable(raw, "description")
			if err != nil {
				return nil, err
			}
		}
		if raw, ok := row["image_url"]; ok && raw != nil {
			imageURL, err = convert.ToStringNullable(raw, "image_url")
			if err != nil {
				return nil, err
			}
		}
		createdUnix, err := convert.ToUint64(createdRaw, "created_at")
		if err != nil {
			return nil, err
		}
		updatedUnix, err := convert.ToUint64(updatedRaw, "updated_at")
		if err != nil {
			return nil, err
		}
		createdAt = time.Unix(int64(createdUnix), 0)
		updatedAt = time.Unix(int64(updatedUnix), 0)
	default:
		return nil, fmt.Errorf("invalid category payload type: %T", resp[0])
	}

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
	return nil, nil
}

func (r *tarantoolRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}
