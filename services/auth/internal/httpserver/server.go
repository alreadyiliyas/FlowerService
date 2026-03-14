package httpserver

import (
	"context"
	"net/http"

	authusecase "github.com/ilyas/flower/services/auth/internal/usecase/auth"
	userusecase "github.com/ilyas/flower/services/auth/internal/usecase/user"
)

// Config описывает настройки HTTP-сервера.
type Config struct {
	Address string
	JWTSecret string
}

// Server инкапсулирует http.Server.
type Server struct {
	httpServer *http.Server
}

// New создаёт новый HTTP-сервер с переданным конфигом.
func New(cfg Config, authUC authusecase.UsecaseAuth, userUC userusecase.UsecaseUser) *Server {
	handler := newRouter(authUC, userUC, cfg.JWTSecret)

	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.Address,
			Handler: handler,
		},
	}
}

// Start запускает HTTP-сервер и слушает до остановки.
func (s *Server) Start(ctx context.Context) error {
	// Грейсфул-шатдаун по завершению контекста.
	go func() {
		<-ctx.Done()
		_ = s.httpServer.Shutdown(context.Background())
	}()

	return s.httpServer.ListenAndServe()
}
