package authcontext

import (
	"context"

	"google.golang.org/grpc"
)

const (
	// ServiceName и GetUserContextPath вынесены в общий пакет, чтобы
	// сервер auth и клиент catalog использовали один и тот же контракт.
	ServiceName        = "auth.v1.AuthService"
	GetUserContextPath = "/auth.v1.AuthService/GetUserContext"
)

// GetUserContextRequest передает в auth токен и, при необходимости,
// session_id из заголовка X-Session-Id для дополнительной сверки.
type GetUserContextRequest struct {
	AccessToken string `json:"access_token"`
	SessionID   string `json:"session_id,omitempty"`
}

// GetUserContextResponse возвращает подтвержденный контекст пользователя,
// который потом кладется в HTTP context в catalog middleware.
type GetUserContextResponse struct {
	UserID          uint64 `json:"user_id"`
	Role            string `json:"role"`
	PhoneNumber     string `json:"phone_number"`
	SessionID       string `json:"session_id"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	IsAuthenticated bool   `json:"is_authenticated"`
}

type AuthServiceServer interface {
	GetUserContext(context.Context, *GetUserContextRequest) (*GetUserContextResponse, error)
}

type AuthServiceClient interface {
	GetUserContext(ctx context.Context, in *GetUserContextRequest, opts ...grpc.CallOption) (*GetUserContextResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc: cc}
}

func (c *authServiceClient) GetUserContext(ctx context.Context, in *GetUserContextRequest, opts ...grpc.CallOption) (*GetUserContextResponse, error) {
	out := new(GetUserContextResponse)
	if err := c.cc.Invoke(ctx, GetUserContextPath, in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterAuthServiceServer регистрирует сервис вручную через ServiceDesc.
// Здесь нет сгенерированного pb-кода, поэтому контракт описан явно.
func RegisterAuthServiceServer(registrar grpc.ServiceRegistrar, srv AuthServiceServer) {
	registrar.RegisterService(&grpc.ServiceDesc{
		ServiceName: ServiceName,
		HandlerType: (*AuthServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetUserContext",
				Handler:    getUserContextHandler,
			},
		},
	}, srv)
}

// getUserContextHandler связывает низкоуровневый gRPC вызов с нашим
// интерфейсом AuthServiceServer и поддерживает interceptor-цепочку.
func getUserContextHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserContextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		return srv.(AuthServiceServer).GetUserContext(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GetUserContextPath,
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).GetUserContext(ctx, req.(*GetUserContextRequest))
	}

	return interceptor(ctx, in, info, handler)
}
