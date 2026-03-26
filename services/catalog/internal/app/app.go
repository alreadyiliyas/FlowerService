package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyas/flower/services/catalog/internal/config"
	authclient "github.com/ilyas/flower/services/catalog/internal/grpc/authclient"
	"github.com/ilyas/flower/services/catalog/internal/httpserver"
	redisclient "github.com/ilyas/flower/services/catalog/internal/redis"
	cacherepo "github.com/ilyas/flower/services/catalog/internal/repositories/cache"
	categoriesrepo "github.com/ilyas/flower/services/catalog/internal/repositories/categories"
	productsrepo "github.com/ilyas/flower/services/catalog/internal/repositories/products"
	tntclient "github.com/ilyas/flower/services/catalog/internal/tarantool"
	usecaseCateg "github.com/ilyas/flower/services/catalog/internal/usecase/categories"
	usecaseProd "github.com/ilyas/flower/services/catalog/internal/usecase/products"
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

	authGRPCClient, err := authclient.New(ctx, cfg.AuthGRPC.Addr)
	if err != nil {
		return fmt.Errorf("failed to connect to auth grpc: %w", err)
	}
	defer authGRPCClient.Close()

	categoriesRepository := categoriesrepo.NewTarantoolRepository(tntConn)
	productsRepository := productsrepo.NewTarantoolRepository(tntConn)
	cacheRepository := cacherepo.NewRedisRepository(redisConn)

	cu := usecaseCateg.NewCategoriesUsecase(categoriesRepository, cacheRepository)
	pu := usecaseProd.NewproductsUsecase(productsRepository, cacheRepository)

	httpSrv := httpserver.New(httpserver.Config{
		Address: cfg.HTTP.Address,
	}, cu, pu, authGRPCClient)

	fmt.Fprintf(os.Stdout, "catalog service listening on %s\n", cfg.HTTP.Address)

	if err := httpSrv.Start(ctx); err != nil {
		return fmt.Errorf("http server stopped: %w", err)
	}

	return nil
}
