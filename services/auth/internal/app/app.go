package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyas/flower/services/auth/internal/config"
	"github.com/ilyas/flower/services/auth/internal/httpserver"
	repositories "github.com/ilyas/flower/services/auth/internal/repositories/auth"
	tntclient "github.com/ilyas/flower/services/auth/internal/tarantool"
	authusecase "github.com/ilyas/flower/services/auth/internal/usecase/auth"
)

func Run() error {
	// Базовый контекст с отменой по сигналу.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	tntConn, err := tntclient.New(ctx, tntclient.Config{
		Addr:     cfg.Tarantool.Addr,
		User:     cfg.Tarantool.User,
		Password: cfg.Tarantool.Password,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to Tarantool: %w", err)
	}
	defer tntConn.Close()

	authRepo := repositories.NewTarantoolRepository(tntConn)
	authUC := authusecase.New(cfg, authRepo)

	httpSrv := httpserver.New(httpserver.Config{
		Address: cfg.HTTP.Address,
	}, authUC)

	fmt.Fprintf(os.Stdout, "auth service listening on %s\n", cfg.HTTP.Address)

	if err := httpSrv.Start(ctx); err != nil {
		return fmt.Errorf("http server stopped: %w", err)
	}

	return nil
}
