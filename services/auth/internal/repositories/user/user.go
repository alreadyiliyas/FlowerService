package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ilyas/flower/services/auth/internal/apperrors"
	"github.com/ilyas/flower/services/auth/internal/entities"
	"github.com/ilyas/flower/services/auth/internal/utils"
	convert "github.com/ilyas/flower/services/auth/internal/utils"
	tnt "github.com/tarantool/go-tarantool/v2"
)

type tarantoolRepository struct {
	conn *tnt.Connection
}

func NewTarantoolRepository(conn *tnt.Connection) UserRepository {
	return &tarantoolRepository{conn: conn}
}

func (r *tarantoolRepository) Get(ctx context.Context, account *entities.Auth) (*entities.User, error) {
	req := tnt.NewCallRequest("get_user_info_by_phone_number").
		Args([]interface{}{
			*account.PhoneNumber,
		}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from get_user_info_by_phone_number")
	}

	row, ok := resp[0].([]interface{})
	if !ok || len(row) < 9 {
		return nil, errors.New("invalid response payload")
	}

	id, err := convert.ToUint64(row[0], "id")
	if err != nil {
		return nil, err
	}
	firstName, err := convert.ToString(row[1], "first_name")
	if err != nil {
		return nil, err
	}
	lastName, err := convert.ToString(row[2], "last_name")
	if err != nil {
		return nil, err
	}
	roleName, err := convert.ToString(row[3], "role")
	if err != nil {
		return nil, err
	}
	isActive, err := convert.ToBool(row[4], "is_active")
	if err != nil {
		return nil, err
	}
	createdUnix, err := convert.ToUint64(row[5], "created_at")
	if err != nil {
		return nil, err
	}
	createdAt := time.Unix(int64(createdUnix), 0)

	updatedUnix, err := convert.ToUint64(row[6], "updated_at")
	if err != nil {
		return nil, err
	}
	updatedAt := time.Unix(int64(updatedUnix), 0)

	avatarURL, err := convert.ToStringNullable(row[7], "avatar_url")
	if err != nil {
		return nil, err
	}
	phoneNumber, err := convert.ToString(row[8], "phone_number")
	if err != nil {
		return nil, err
	}

	return &entities.User{
		Id:          &id,
		FirstName:   &firstName,
		LastName:    &lastName,
		PhoneNumber: &phoneNumber,
		Role:        &roleName,
		AvatarURL:   avatarURL,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
		IsActive:    isActive,
	}, nil
}

func (r *tarantoolRepository) Update(ctx context.Context, account *entities.Auth) (*entities.User, error) {
	req := tnt.NewCallRequest("update_user_info_by_phone_number").
		Args([]interface{}{
			*account.PhoneNumber,
			utils.ValueOrNull(account.User.PhoneNumber),
			*account.User.FirstName,
			*account.User.LastName,
			utils.ValueOrNull(account.User.AvatarURL),
		}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from update_user_info_by_phone_number")
	}

	row, ok := resp[0].([]interface{})
	if !ok || len(row) < 9 {
		return nil, errors.New("invalid response payload")
	}

	id, err := convert.ToUint64(row[0], "id")
	if err != nil {
		return nil, err
	}
	firstName, err := convert.ToString(row[1], "first_name")
	if err != nil {
		return nil, err
	}
	lastName, err := convert.ToString(row[2], "last_name")
	if err != nil {
		return nil, err
	}
	roleName, err := convert.ToString(row[3], "role")
	if err != nil {
		return nil, err
	}
	isActive, err := convert.ToBool(row[4], "is_active")
	if err != nil {
		return nil, err
	}
	createdUnix, err := convert.ToUint64(row[5], "created_at")
	if err != nil {
		return nil, err
	}
	createdAt := time.Unix(int64(createdUnix), 0)

	updatedUnix, err := convert.ToUint64(row[6], "updated_at")
	if err != nil {
		return nil, err
	}
	updatedAt := time.Unix(int64(updatedUnix), 0)

	avatarURL, err := convert.ToStringNullable(row[7], "avatar_url")
	if err != nil {
		return nil, err
	}
	phoneNumber, err := convert.ToString(row[8], "phone_number")
	if err != nil {
		return nil, err
	}

	return &entities.User{
		Id:          &id,
		FirstName:   &firstName,
		LastName:    &lastName,
		PhoneNumber: &phoneNumber,
		Role:        &roleName,
		AvatarURL:   avatarURL,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
		IsActive:    isActive,
	}, nil
}

func (r *tarantoolRepository) Delete(ctx context.Context, account *entities.Auth) error {
	return r.callBoolProc(ctx, "delete_user", []interface{}{*account.PhoneNumber})
}

func (r *tarantoolRepository) callBoolProc(ctx context.Context, proc string, args []interface{}) error {
	req := tnt.NewCallRequest(proc).Args(args).Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return fmt.Errorf("%w: %s returned empty response", apperrors.ErrDB, proc)
	}

	ok, isBool := resp[0].(bool)
	if !isBool {
		return fmt.Errorf("%w: invalid %s response type %T", apperrors.ErrDB, proc, resp[0])
	}
	if !ok {
		return fmt.Errorf("%w: %s returned false", apperrors.ErrDB, proc)
	}

	return nil
}
