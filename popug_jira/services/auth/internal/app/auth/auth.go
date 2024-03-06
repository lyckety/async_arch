package producer

import (
	"context"

	"github.com/lyckety/async_arch/popug_jira/services/auth/internal/db/domain"
	pbV1 "github.com/lyckety/async_arch/popug_jira/services/auth/pkg/grpc/auth/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	pbV1.UnimplementedAuthServiceServer

	dbIns      domain.Repository
	jwtManager JWTManager
}

func New(db domain.Repository, jwt JWTManager) *AuthService {
	return &AuthService{
		dbIns:      db,
		jwtManager: jwt,
	}
}

func (a *AuthService) UserLogin(
	ctx context.Context,
	req *pbV1.UserLoginRequest,
) (*pbV1.UserLoginResponse, error) {
	userInfo, err := a.dbIns.GetUserByUsername(ctx, req.Username)
	if err != nil {
		log.Errorf("a.dbIns.GetUserByUsername(ctx, %s): %s", req.Username, err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"a.dbIns.GetUserByUsername(ctx, %s): %s", req.Username, err,
		)
	}

	if userInfo.Password != req.Password {
		return nil, status.Error(
			codes.Unauthenticated,
			"incorrect username or password",
		)
	}

	tokenString, err := a.jwtManager.GenerateToken(req.Username, userInfo.PublicID.String(), string(userInfo.Role))
	if err != nil {
		log.Errorf("jwt.NewWithClaims(...): %s", err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"jwt.NewWithClaims(...): %s", err.Error(),
		)
	}

	return &pbV1.UserLoginResponse{
		Token: tokenString,
	}, nil
}

type JWTManager interface {
	GenerateToken(username, publicID, role string) (string, error)
}
