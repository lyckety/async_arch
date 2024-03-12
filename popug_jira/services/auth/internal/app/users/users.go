package producer

import (
	"context"
	"errors"

	events "github.com/lyckety/async_arch/popug_jira/services/auth/internal/app/users/events"
	"github.com/lyckety/async_arch/popug_jira/services/auth/internal/db/domain"
	pbV1User "github.com/lyckety/async_arch/popug_jira/services/auth/pkg/api/grpc/user/v1"
	pbV1Users "github.com/lyckety/async_arch/popug_jira/services/auth/pkg/api/grpc/users/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrUserRoleMustBeSet = errors.New("user role must be set")

type UsersService struct {
	pbV1Users.UnimplementedUsersServiceServer

	dbIns          domain.Repository
	producerEvents *events.UserCUDEventSender
}

func New(db domain.Repository, cudEventSender *events.UserCUDEventSender) *UsersService {
	return &UsersService{
		dbIns:          db,
		producerEvents: cudEventSender,
	}
}

func (s *UsersService) CreateUser(
	ctx context.Context,
	req *pbV1Users.CreateUserRequest,
) (*pbV1Users.CreateUserResponse, error) {
	userDB := &domain.User{
		ID:        [16]byte{},
		PublicID:  [16]byte{},
		FirstName: req.User.GetFirstName(),
		LastName:  req.User.GetLastName(),
		Username:  req.User.GetUsername(),
		Password:  req.GetPassword(),
		Email:     req.User.GetEmail(),
	}

	userRoleDB := s.convertRPCRoleToDBRole(req.User.GetRole())
	if userRoleDB == domain.Unknown {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"s.convertRPCRoleToDBRole(%s): %s", req.User.GetRole().String(), ErrUserRoleMustBeSet,
		)
	}

	userDB.Role = userRoleDB

	createdUser, err := s.dbIns.CreateUser(ctx, userDB)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"s.dbIns.CreateUser(ctx, %v): %s", userDB, err.Error(),
		)
	}

	eventInfo := events.CreatedUserV1{
		EventTime: createdUser.CreatedAt.Unix(),
		PublicID:  createdUser.PublicID,
		Username:  createdUser.Username,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Role:      string(createdUser.Role),
		Email:     createdUser.Email,
	}

	go func() {
		if err := s.producerEvents.Send(context.Background(), eventInfo); err != nil {
			log.Errorf("failed send event for create or update user %v: %s", req.User, err.Error())

			return
		}

		log.Debugf("success sent cud event (user created): %v", eventInfo)
	}()

	return &pbV1Users.CreateUserResponse{Id: createdUser.PublicID.String()}, nil
}

func (s *UsersService) UpdateUser(
	ctx context.Context,
	req *pbV1Users.UpdateUserRequest,
) (*pbV1Users.UpdateUserResponse, error) {
	userDB := &domain.User{
		ID:        [16]byte{},
		PublicID:  [16]byte{},
		FirstName: req.User.GetFirstName(),
		LastName:  req.User.GetLastName(),
		Username:  req.User.GetUsername(),
		Email:     req.User.GetEmail(),
	}

	userRoleDB := s.convertRPCRoleToDBRole(req.User.GetRole())
	if userRoleDB == domain.Unknown {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"s.convertRPCRoleToDBRole(%s): %s", req.User.GetRole().String(), ErrUserRoleMustBeSet,
		)
	}

	userDB.Role = userRoleDB

	updatedUser, err := s.dbIns.UpdateUser(ctx, userDB)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"s.dbIns.UpdateUser(ctx, %v): %s", userDB, err.Error(),
		)
	}

	eventInfo := events.UpdatedUserV1{
		EventTime: updatedUser.CreatedAt.Unix(),
		PublicID:  updatedUser.PublicID,
		Username:  updatedUser.Username,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Role:      string(updatedUser.Role),
		Email:     updatedUser.Email,
	}

	go func() {
		if err := s.producerEvents.Send(context.Background(), eventInfo); err != nil {
			log.Errorf("failed send event for create or update user %v: %s", req.User, err.Error())

			return
		}

		log.Debugf("success sent cud event (user updated): %v", eventInfo)
	}()

	return &pbV1Users.UpdateUserResponse{Id: updatedUser.PublicID.String()}, nil
}

func (s *UsersService) convertRPCRoleToDBRole(rpcRole pbV1User.UserRole) domain.UserRoleType {
	switch rpcRole {
	case pbV1User.UserRole_USER_ROLE_ADMIN:
		return domain.Administrator
	case pbV1User.UserRole_USER_ROLE_BOOKKEEPER:
		return domain.Bookkeepper
	case pbV1User.UserRole_USER_ROLE_MANAGER:
		return domain.Manager
	case pbV1User.UserRole_USER_ROLE_WORKER:
		return domain.Worker
	case pbV1User.UserRole_USER_ROLE_UNSPECIFIED:
		fallthrough
	default:
		return domain.Unknown
	}
}
