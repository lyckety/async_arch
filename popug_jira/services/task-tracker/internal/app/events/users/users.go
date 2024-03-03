package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/db/domain"
	mbCons "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/mb/consumer"
	pbV1Events "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/grpc/authevents/v1"
	pbV1User "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/grpc/user/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type UserCUDEventProcessor struct {
	mbClient *mbCons.MBConsumer
	storage  domain.UserRepository
}

func New(
	repository domain.UserRepository,
	brokerURLs []string,
	topic string,
	groupID string,
	partition int,
) *UserCUDEventProcessor {
	return &UserCUDEventProcessor{
		mbClient: mbCons.New(brokerURLs, topic, groupID, partition),
		storage:  repository,
	}
}

func (p *UserCUDEventProcessor) StartProcessMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		message, err := p.mbClient.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, mbCons.ErrContextCanceled) {
				logrus.Debug("p.mbClient.FetchMessage(ctx): context canceled")
			} else {
				logrus.Errorf("UserCUDEventProcessor: error fetch message from mb: %s", err)
			}

			continue
		}

		logrus.Debugf("UserCUDEventProcessor: success fetch message %v from mb: %s", message, err)

		rpcAuthEvent := &pbV1Events.AuthEvent{}

		if err := proto.Unmarshal(message.Value, rpcAuthEvent); err != nil {
			logrus.Errorf("UserCUDEventProcessor: error unmarshal message from mb to protobuf: %s", err)

			continue
		}

		switch rpcAuthEvent.GetEventType() {
		case pbV1Events.AuthCUDEventType_AUTH_CUD_EVENT_TYPE_CREATE,
			pbV1Events.AuthCUDEventType_AUTH_CUD_EVENT_TYPE_UPDATE:
			userDB, err := convertRPCUserToDB(rpcAuthEvent.GetUser())
			if err != nil {
				logrus.Errorf("error convert rpc user to database user: %s", err)

				continue
			}

			if _, err := p.storage.CreateOrUpdateUser(ctx, userDB); err != nil {
				logrus.Errorf("error create or update user in database: %s", err)

				continue
			}
		default:
			// Очень плохо...
			logrus.Fatalf("UserCUDEventProcessor: unsupported event type: %s", err)
		}

		if err := p.mbClient.CommitMessages(ctx, message); err != nil {
			logrus.Errorf("UserCUDEventProcessor: error commit message %v from mb: %s", message, err)

			continue
		}

		logrus.Debugf("UserCUDEventProcessor: success commit message %v from mb!", message)

		logrus.Debugf("UserCUDEventProcessor: success processing message %v from mb!", message)
	}
}

func convertRPCUserToDB(userWithID *pbV1User.UserWithID) (*domain.User, error) {
	uuidDB, err := uuid.Parse(userWithID.GetId())
	if err != nil {
		return nil, fmt.Errorf("rpc user must have id uuid type: %w", err)
	}

	userRoleDB := convertRPCUserRoleToDB(userWithID.GetUser().GetRole())
	if userRoleDB == domain.UnknownRole {
		return nil, fmt.Errorf("rpc user have unknown type user role")
	}

	return &domain.User{
		ID:        uuidDB,
		FirstName: userWithID.GetUser().GetFirstName(),
		LastName:  userWithID.GetUser().GetLastName(),
		Username:  userWithID.GetUser().GetUsername(),
		Email:     userWithID.GetUser().GetEmail(),
		Role:      userRoleDB,
	}, nil

}

func convertRPCUserRoleToDB(rpcRole pbV1User.UserRole) domain.UserRoleType {
	switch rpcRole {
	case pbV1User.UserRole_USER_ROLE_ADMIN:
		return domain.AdministratorRole
	case pbV1User.UserRole_USER_ROLE_BOOKKEEPER:
		return domain.BookkeepperRole
	case pbV1User.UserRole_USER_ROLE_MANAGER:
		return domain.ManagerRole
	case pbV1User.UserRole_USER_ROLE_WORKER:
		return domain.WorkerRole
	case pbV1User.UserRole_USER_ROLE_UNSPECIFIED:
		fallthrough
	default:
		return domain.UnknownRole
	}
}
