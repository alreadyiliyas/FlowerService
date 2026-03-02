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
	// TODO: Реализовать получение учетной записи по номеру телефона в Тарантуле
	return nil, errors.New("not implemented")
}

func (r *tarantoolRepository) VerifyAccount(ctx context.Context, phoneNumber *string) error {
	req := tnt.NewCallRequest("verify_account").
		Args([]interface{}{*phoneNumber}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return convert.MapTarantoolError(err)
	}
	if len(resp) == 0 {
		return fmt.Errorf("%w: verify_account returned empty response", apperrors.ErrDB)
	}

	ok, isBool := resp[0].(bool)
	if !isBool {
		return fmt.Errorf("%w: invalid verify_account response type %T", apperrors.ErrDB, resp[0])
	}
	if !ok {
		return fmt.Errorf("%w: user didn't activate", apperrors.ErrDB)
	}

	return nil
}

func (r *tarantoolRepository) SetPassword(ctx context.Context, account *entities.Auth) error {
	req := tnt.NewCallRequest("set_password").
		Args([]interface{}{*account.PhoneNumber, *account.PasswordHash}).
		Context(ctx)

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return convert.MapTarantoolError(err)
	}

	if len(resp) == 0 {
		return fmt.Errorf("%w: set_password returned empty response", apperrors.ErrDB)
	}

	ok, isBool := resp[0].(bool)
	if !isBool {
		return fmt.Errorf("%w: invalid set_password response type %T", apperrors.ErrDB, resp[0])
	}
	if !ok {
		return fmt.Errorf("%w: password of user didn't set", apperrors.ErrDB)
	}

	return nil
}

func (r *tarantoolRepository) GetPassword(ctx context.Context, in *string) (*string, error) {
	// TODO: Реализовать получение пароля из Тарантула
	return nil, errors.New("not implemented")
}
