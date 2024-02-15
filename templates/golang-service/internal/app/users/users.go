package producer

import (
	"context"

	"github.com/lyckety/async_arch/golang-service/internal/db/domain"
	pbV1 "github.com/lyckety/async_arch/golang-service/pkg/grpc/users/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UsersService struct {
	pbV1.UnimplementedUsersServiceServer

	dbIns domain.Repository
}

func New(db domain.Repository) *UsersService {
	return &UsersService{
		dbIns: db,
	}
}

func (s *UsersService) CreateOrUpdateUser(
	ctx context.Context,
	req *pbV1.CreateOrUpdateUserRequest,
) (*pbV1.CreateOrUpdateUserResponse, error) {
	userDB := &domain.User{
		Username:    req.GetUser().GetUsername(),
		Email:       req.GetUser().GetEmail(),
		DateOfBirth: req.GetUser().GetDateOfBirth().AsTime(),
	}

	if err := s.dbIns.CreateOrUpdateUser(ctx, userDB); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"s.dbIns.CreateOrUpdateUser(ctx, %v): %s", userDB, err.Error(),
		)
	}

	return &pbV1.CreateOrUpdateUserResponse{}, nil
}

func (s *UsersService) DeleteUser(
	ctx context.Context,
	req *pbV1.DeleteUserRequest,
) (*pbV1.DeleteUserResponse, error) {
	if req.GetUsername() == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"username must be set",
		)
	}

	if err := s.dbIns.DeleteUser(ctx, req.GetUsername()); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"s.dbIns.DeleteUser(ctx, %v): %s", req.GetUsername(), err.Error(),
		)
	}

	return &pbV1.DeleteUserResponse{}, nil
}

func (s *UsersService) GetAllUsers(
	ctx context.Context,
	req *pbV1.GetAllUsersRequest,
) (*pbV1.GetAllUsersResponse, error) {
	dbUsers, err := s.dbIns.GetAllUsers(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"s.dbIns.GetAllUsers(ctx): %s", err.Error(),
		)
	}

	resp := &pbV1.GetAllUsersResponse{
		Users: make([]*pbV1.User, len(dbUsers)),
	}

	for i, dbUser := range dbUsers {
		resp.Users[i] = &pbV1.User{
			Username: dbUser.Username,
			Email:    dbUser.Email,
			DateOfBirth: &timestamppb.Timestamp{
				Seconds: dbUser.DateOfBirth.Unix(),
				Nanos:   0,
			},
		}
	}

	return resp, nil
}

func (s *UsersService) GetUserByName(
	ctx context.Context,
	req *pbV1.GetUserByNameRequest,
) (*pbV1.GetUserByNameResponse, error) {
	if req.GetUsername() == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"username must be set",
		)
	}

	dbUser, err := s.dbIns.GetUserByName(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"s.dbIns.GetUserByName(ctx, %s): %s", req.GetUsername(), err.Error(),
		)
	}

	userRPC := &pbV1.User{
		Username: dbUser.Username,
		Email:    dbUser.Email,
		DateOfBirth: &timestamppb.Timestamp{
			Seconds: dbUser.DateOfBirth.Unix(),
			Nanos:   0,
		},
	}

	return &pbV1.GetUserByNameResponse{
		User: userRPC,
	}, nil
}
