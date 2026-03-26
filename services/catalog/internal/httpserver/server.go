package httpserver

import (
	"context"
	"net/http"

	authclient "github.com/ilyas/flower/services/catalog/internal/grpc/authclient"
	usecaseCateg "github.com/ilyas/flower/services/catalog/internal/usecase/categories"
	usecaseProd "github.com/ilyas/flower/services/catalog/internal/usecase/products"
)

// Config описывает настройки HTTP-сервера.
type Config struct {
	Address string
}

// Server инкапсулирует http.Server.
type Server struct {
	httpServer *http.Server
}

// New создаёт новый HTTP-сервер с переданным конфигом.
func New(cfg Config, cu usecaseCateg.UsecaseCategories, pu usecaseProd.ProductUsecase, authClient authclient.Client) *Server {
	handler := newRouter(cu, pu, authClient)

	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.Address,
			Handler: handler,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		_ = s.httpServer.Shutdown(context.Background())
	}()

	return s.httpServer.ListenAndServe()
}
