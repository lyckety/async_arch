package interceptors

import (
	"context"

	"github.com/lyckety/async_arch/popug_jira/services/auth/internal/db/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func JWTUnary(
	jwtMgr JWTManager,
) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	jwtAuth := jwtMgr

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/auth.v1.AuthService/Login" {
			return handler(ctx, req)
		}

		role, err := jwtAuth.Authorize(ctx)
		if err != nil {
			return nil, err
		}

		switch info.FullMethod {
		case "/users.v1.UsersService/CreateUser",
			"/users.v1.UsersService/UpdateUser":
			if role != string(domain.Administrator) {
				return nil, status.Error(codes.PermissionDenied, "no permission to access this method")
			}
		default:
			return nil, status.Errorf(
				codes.PermissionDenied,
				"no permission to access this method: unknown api method %q", info.FullMethod,
			)
		}

		return handler(ctx, req)
	}
}

type JWTManager interface {
	Authorize(ctx context.Context) (string, error)
}
