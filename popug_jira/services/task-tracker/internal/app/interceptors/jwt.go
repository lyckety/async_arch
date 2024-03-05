package interceptors

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/db/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	jwtSecretKey     = "09876543210"
	ContextKeyUserID = "user_id"
)

type ExtendedClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	ID       string `json:"id"`
	Role     string `json:"role"`
}

func JWTUnary(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	claims, err := validateJWTToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("error validate jwt token: %w", err)
	}

	switch info.FullMethod {
	case "/tasktracker.v1.TaskTrackerService/CreateTask":
		// Новые таски может создавать кто угодно (администратор, начальник, разработчик, менеджер и любая другая роль).
	case "/tasktracker.v1.TaskTrackerService/RandomReassignOpenedTasks":
		if claims.Role != string(domain.AdministratorRole) && claims.Role != string(domain.ManagerRole) {
			return nil, fmt.Errorf("unathorized access")
		}
	case "/tasktracker.v1.TaskTrackerService/TaskComplete",
		"/tasktracker.v1.TaskTrackerService/GetListOpenedTasksForMe":
		if claims.Role != string(domain.WorkerRole) {
			return nil, fmt.Errorf("unathorized access")
		}
	}

	userID := claims.ID

	ctx = context.WithValue(ctx, ContextKeyUserID, userID)

	return handler(ctx, req)
}

func validateJWTToken(ctx context.Context) (ExtendedClaims, error) {
	secretKey := []byte(jwtSecretKey)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ExtendedClaims{}, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return ExtendedClaims{}, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	// The authHeader should be in the format `Bearer <token>`
	splitToken := strings.Split(authHeader[0], " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return ExtendedClaims{}, status.Errorf(codes.Unauthenticated, "authorization token format is invalid")
	}

	tokenStr := splitToken[1]

	claims := &ExtendedClaims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return ExtendedClaims{}, fmt.Errorf("invalid token signature: %w", err)
		}

		return ExtendedClaims{}, fmt.Errorf("invalid token: %w", err)
	}

	if !tkn.Valid {
		return ExtendedClaims{}, fmt.Errorf("invalid token: %w", err)
	}

	return *claims, nil
}
