package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	tnt "github.com/tarantool/go-tarantool/v2"

	"github.com/ilyas/flower/services/auth/internal/apperrors"
	"github.com/ilyas/flower/services/auth/internal/entities"
	convert "github.com/ilyas/flower/services/auth/internal/utils"
)

type tarantoolRepository struct {
	conn *tnt.Connection
}

func NewTarantoolRepository(conn *tnt.Connection) AuthRepository {
	return &tarantoolRepository{conn: conn}
}

func (r *tarantoolRepository) CreateUser(ctx context.Context, user *entities.User, account *entities.Auth) (*entities.User, error) {
	req := tnt.NewCallRequest("create_user").
		Args([]interface{}{
			*user.FirstName,
			*user.LastName,
			*account.PhoneNumber,
			*user.Role,
		}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from create_user")
	}

	row, ok := resp[0].([]interface{})
	if !ok || len(row) < 5 {
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
	createdUnix, err := convert.ToUint64(row[4], "created_at")
	if err != nil {
		return nil, err
	}
	createdAt := time.Unix(int64(createdUnix), 0)

	account.UserId = &id

	return &entities.User{
		Id:        &id,
		FirstName: &firstName,
		LastName:  &lastName,
		Role:      &roleName,
		CreatedAt: &createdAt,
		IsActive:  false,
		Version:   1,
	}, nil
}

func (r *tarantoolRepository) GetAccountByPhoneNumber(ctx context.Context, in *string) (*entities.Auth, error) {
	req := tnt.NewCallRequest("get_account_by_phone_number").
		Args([]interface{}{
			*in,
		}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from get account by phone number")
	}

	row, ok := resp[0].([]interface{})
	if !ok || len(row) < 5 {
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
	passwordHash, err := convert.ToString(row[4], "password_hash")
	if err != nil {
		return nil, err
	}
	user := entities.User{
		Id:        &id,
		FirstName: &firstName,
		LastName:  &lastName,
		Role:      &roleName,
	}

	return &entities.Auth{
		UserId:       &id,
		PhoneNumber:  in,
		PasswordHash: &passwordHash,
		User:         user,
	}, nil
}

func (r *tarantoolRepository) VerifyAccount(ctx context.Context, phoneNumber *string) error {
	return r.callBoolProc(ctx, "verify_account", []interface{}{*phoneNumber})
}

func (r *tarantoolRepository) SetPassword(ctx context.Context, account *entities.Auth) error {
	return r.callBoolProc(ctx, "set_password", []interface{}{*account.PhoneNumber, *account.PasswordHash})
}

func (r *tarantoolRepository) UpdatePassword(ctx context.Context, account *entities.Auth) error {
	return r.callBoolProc(ctx, "update_password", []interface{}{*account.PhoneNumber, *account.PasswordHash})
}

func (r *tarantoolRepository) GetPassword(ctx context.Context, in *string) (*string, error) {
	return nil, errors.New("not implemented")
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
