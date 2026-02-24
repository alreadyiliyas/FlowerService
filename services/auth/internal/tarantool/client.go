package tarantool

import (
	"context"
	"time"

	tnt "github.com/tarantool/go-tarantool/v2"
)

// Config описывает параметры подключения к Tarantool.
type Config struct {
	Addr     string
	User     string
	Password string
}

// New создаёт новое подключение к Tarantool.
func New(ctx context.Context, cfg Config) (*tnt.Connection, error) {
	// ToDo вынести в конфиг
	opts := tnt.Opts{
		Timeout:       3 * time.Second,
		Reconnect:     1 * time.Second,
		MaxReconnects: 5,
	}
	dialer := tnt.NetDialer{
		Address:  cfg.Addr,
		User:     cfg.User,
		Password: cfg.Password,
	}

	conn, err := tnt.Connect(ctx, dialer, opts)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
