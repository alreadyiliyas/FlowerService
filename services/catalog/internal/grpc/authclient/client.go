package authclient

import (
	"context"
	"time"

	grpcauth "github.com/ilyas/flower/pkg/grpc/authcontext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	GetUserContext(ctx context.Context, accessToken, sessionID string) (*grpcauth.GetUserContextResponse, error)
	Close() error
}

type grpcClient struct {
	conn   *grpc.ClientConn
	client grpcauth.AuthServiceClient
}

// codec из pkg/grpc/authcontext чтобы одинаково сериализовал запросы без генерации pb на этом этапе.
func New(ctx context.Context, addr string) (Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(grpcauth.Codec()), grpc.WaitForReady(true)),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClient{
		conn:   conn,
		client: grpcauth.NewAuthServiceClient(conn),
	}, nil
}

// Получает подтвержденный user context для middleware.
func (c *grpcClient) GetUserContext(ctx context.Context, accessToken, sessionID string) (*grpcauth.GetUserContextResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.GetUserContext(callCtx, &grpcauth.GetUserContextRequest{
		AccessToken: accessToken,
		SessionID:   sessionID,
	}, grpc.ForceCodec(grpcauth.Codec()))
}

// Close освобождает underlying grpc connection.
func (c *grpcClient) Close() error {
	return c.conn.Close()
}
