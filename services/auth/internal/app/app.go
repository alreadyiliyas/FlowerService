package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyas/flower/services/auth/internal/config"
	"github.com/ilyas/flower/services/auth/internal/httpserver"
	redisclient "github.com/ilyas/flower/services/auth/internal/redis"
	authRepo "github.com/ilyas/flower/services/auth/internal/repositories/auth"
	cacheRepo "github.com/ilyas/flower/services/auth/internal/repositories/cache"
	userRepo "github.com/ilyas/flower/services/auth/internal/repositories/user"
	tntclient "github.com/ilyas/flower/services/auth/internal/tarantool"
	authusecase "github.com/ilyas/flower/services/auth/internal/usecase/auth"
	userusecase "github.com/ilyas/flower/services/auth/internal/usecase/user"
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

	redisConn, err := redisclient.New(ctx, redisclient.Config{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	defer redisConn.Close()

	authRepo := authRepo.NewTarantoolRepository(tntConn)
	userRepo := userRepo.NewTarantoolRepository(tntConn)
	cacheRepo := cacheRepo.NewRedisRepository(redisConn)
	authUC := authusecase.New(cfg, authRepo, cacheRepo)
	userUC := userusecase.New(cfg, userRepo, cacheRepo)

	httpSrv := httpserver.New(httpserver.Config{
		Address:   cfg.HTTP.Address,
		JWTSecret: cfg.JWT.Secret,
	}, authUC, userUC)

	fmt.Fprintf(os.Stdout, "auth service listening on %s\n", cfg.HTTP.Address)

	if err := httpSrv.Start(ctx); err != nil {
		return fmt.Errorf("http server stopped: %w", err)
	}

	return nil
}
