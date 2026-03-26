package grpcserver

import (
	"context"
	"net"

	grpcauth "github.com/ilyas/flower/pkg/grpc/authcontext"
	authusecase "github.com/ilyas/flower/services/auth/internal/usecase/auth"
	"google.golang.org/grpc"
)

type Config struct {
	Address string
}

// Server инкапсулирует отдельный gRPC сервер auth.
type Server struct {
	address    string
	grpcServer *grpc.Server
}

// New поднимает gRPC server и регистрирует в нем auth сервис.
func New(cfg Config, authUC authusecase.UsecaseAuth) *Server {
	server := grpc.NewServer()
	grpcauth.RegisterAuthServiceServer(server, NewAuthHandler(authUC))

	return &Server{
		address:    cfg.Address,
		grpcServer: server,
	}
}

// Start слушает TCP адрес и завершает gRPC server вместе с общим контекстом приложения.
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		s.grpcServer.GracefulStop()
	}()

	if err := s.grpcServer.Serve(lis); err != nil && ctx.Err() == nil {
		return err
	}
	return nil
}
