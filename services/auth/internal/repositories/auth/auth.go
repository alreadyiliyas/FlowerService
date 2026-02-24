package repositories

import (
	"context"
	"errors"
	"time"

	tnt "github.com/tarantool/go-tarantool/v2"

	"github.com/ilyas/flower/services/auth/internal/entities"
)

type tarantoolRepository struct {
	conn *tnt.Connection
}

func NewTarantoolRepository(conn *tnt.Connection) AuthRepository {
	return &tarantoolRepository{conn: conn}
}

func (r *tarantoolRepository) CreateUser(ctx context.Context, user *entities.User, account *entities.Auth) (*entities.User, error) {
	req := tnt.NewCallRequest("create_user").Args([]interface{}{
		*user.FirstName,
		*user.LastName,
		*account.PhoneNumber,
		*user.Role,
	})

	resp, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, errors.New("empty response from create_user")
	}
	row, ok := resp[0].([]interface{})
	if !ok || len(row) < 5 {
		return nil, errors.New("invalid response payload")
	}

	id := row[0].(uint64)
	firstName := row[1].(string)
	lastName := row[2].(string)
	roleName := row[3].(string)
	createdUnix := row[4].(uint64)
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
	// TODO: реализовать получение учетной записи по номеру телефона в Тарантуле
	return nil, errors.New("not implemented")
}

func (r *tarantoolRepository) VerifyAccount(ctx context.Context, account *entities.Auth) error {
	// TODO: реализовать верификацию аккаунта в Тарантуле
	return errors.New("not implemented")
}

func (r *tarantoolRepository) SetPassword(ctx context.Context, account *entities.Auth) error {
	// TODO: реализовать установку пароля в Тарантуле
	return errors.New("not implemented")
}

func (r *tarantoolRepository) GetPassword(ctx context.Context, in *string) (*string, error) {
	// TODO: реализовать получение пароля из Тарантула
	return nil, errors.New("not implemented")
}
