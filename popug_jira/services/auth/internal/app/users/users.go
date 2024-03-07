package producer

import (
	"context"
	"errors"

	events "github.com/lyckety/async_arch/popug_jira/services/auth/internal/app/users/events"
	"github.com/lyckety/async_arch/popug_jira/services/auth/internal/db/domain"
	pbV1Events "github.com/lyckety/async_arch/popug_jira/services/auth/pkg/grpc/authevents/v1"
	pbV1User "github.com/lyckety/async_arch/popug_jira/services/auth/pkg/grpc/user/v1"
	pbV1Users "github.com/lyckety/async_arch/popug_jira/services/auth/pkg/grpc/users/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrUserRoleMustBeSet = errors.New("user role must be set")

type UsersService struct {
	pbV1Users.UnimplementedUsersServiceServer

	dbIns          domain.Repository
	producerEvents *events.CUDEventSender
}

func New(db domain.Repository, cudEventSender *events.CUDEventSender) *UsersService {
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

	eventsMsg := &pbV1Events.AuthEvent{
		EventType: pbV1Events.AuthCUDEventType_AUTH_CUD_EVENT_TYPE_CREATE,
		User: &pbV1User.UserWithID{
			Id:   createdUser.PublicID.String(),
			User: req.GetUser(),
		},

		Timestamp: userDB.CreatedAt.Unix(),
	}

	go func() {
		if err := s.producerEvents.Send(context.Background(), eventsMsg); err != nil {
			log.Errorf("failed send event for create or update user %v: %s", req.User, err.Error())

			return
		}

		log.Debugf("success sent cud event (user created): %v", eventsMsg)
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

	eventsMsg := &pbV1Events.AuthEvent{
		EventType: pbV1Events.AuthCUDEventType_AUTH_CUD_EVENT_TYPE_UPDATE,
		User: &pbV1User.UserWithID{
			Id:   updatedUser.PublicID.String(),
			User: req.GetUser(),
		},

		Timestamp: userDB.UpdatedAt.Unix(),
	}

	go func() {
		if err := s.producerEvents.Send(context.Background(), eventsMsg); err != nil {
			log.Errorf("failed send event for create or update user %v: %s", req.User, err.Error())

			return
		}

		log.Debugf("success sent cud event (user updated): %v", eventsMsg)
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
