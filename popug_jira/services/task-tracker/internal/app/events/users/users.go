package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/db/domain"
	mbCons "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/mb/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	// pbV1Header "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/eventheaders/header/v1"
	pbV1Created "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/created/v1"
	pbV1Updated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/updated/v1"

	// pbV1Header "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/eventheaders/header/v1"

	pbV1UserRole "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/userrole/v1"
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

				if err := p.mbClient.CommitMessages(ctx, message); err != nil {
					logrus.Errorf(
						"p.mbClient.CommitMessages(...): error commit message %v from mb: %s", message, err.Error(),
					)
				}
			}

			continue
		}

		logrus.Debugf("UserCUDEventProcessor: success fetched message %v from mb", message)

		if err := p.processMessage(ctx, message); err != nil {
			logrus.Errorf(
				"p.processMessage(...): error commit message %v from mb: %s", message, err.Error(),
			)
		}

		logrus.Debugf("UserCUDEventProcessor: success processed message %v from mb!", message)
	}
}

func (p *UserCUDEventProcessor) processMessage(ctx context.Context, msg kafka.Message) error {
	var err error
	var dbUser *domain.User

	switch pbData := getUnmarshalledPbFroBinary(msg.Value).(type) {
	case *pbV1Created.Event:
		dbUser, err = convertCreatedUserToDB(pbData)
		if err != nil {
			return fmt.Errorf("processMessage(msg): %v :%w", pbData, err)
		}
	case *pbV1Updated.Event:
		dbUser, err = convertUpdatedUserToDB(pbData)
		if err != nil {
			return fmt.Errorf("processMessage(msg): %q :%w", pbData, err)
		}
	default:
		return fmt.Errorf("processMessage(msg): unknown protobuf type for event message")
	}

	if _, err := p.storage.CreateOrUpdateUser(ctx, dbUser); err != nil {
		return fmt.Errorf("processMessage(msg): error p.storage.CreateOrUpdateUser(ctx, %v): %w", dbUser, err)
	}

	if err := p.mbClient.CommitMessages(ctx, msg); err != nil {
		return fmt.Errorf("p.mbClient.CommitMessages(...): error commit message %v from mb: %w", msg, err)
	}

	logrus.Debugf("p.mbClient.CommitMessages(...): success commit message %v from mb!", msg)

	return nil
}

func getUnmarshalledPbFroBinary(data []byte) protoreflect.ProtoMessage {
	userCreatedV1 := &pbV1Created.Event{}
	userUpdatedV1 := &pbV1Updated.Event{}

	//
	if err := proto.Unmarshal(data, userCreatedV1); err == nil && userCreatedV1.GetHeader().GetEventName() == "UserCreated" { // TODO: очень глупая ошибка, пока как костыль. Но надо исправить
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			userCreatedV1.GetHeader().GetEventName(),
			userCreatedV1,
		)

		return userCreatedV1
	} else if err := proto.Unmarshal(data, userUpdatedV1); err == nil && userCreatedV1.GetHeader().GetEventName() == "UserUpdated" { // TODO: очень глупая ошибка, пока как костыль. Но надо исправить
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			userUpdatedV1.GetHeader().GetEventName(),
			userUpdatedV1,
		)

		return userUpdatedV1
	} else {
		return nil
	}
}

func convertCreatedUserToDB(pbEvent *pbV1Created.Event) (*domain.User, error) {
	uuidDB, err := uuid.Parse(pbEvent.GetData().GetPublicId())
	if err != nil {
		return nil, fmt.Errorf("rpc user must have id uuid type: %w", err)
	}

	userRoleDB := convertRPCUserRoleToDB(pbEvent.GetData().GetRole())
	if userRoleDB == domain.UnknownRole {
		logrus.Errorf("rpc created user event: user have unknown type user role")

		return nil, fmt.Errorf("rpc created user event: user have unknown type user role")
	}

	return &domain.User{
		ID:        uuidDB,
		FirstName: pbEvent.GetData().GetFirstName(),
		LastName:  pbEvent.GetData().GetLastName(),
		Username:  pbEvent.GetData().GetUsername(),
		Email:     pbEvent.GetData().GetEmail(),
		Role:      userRoleDB,
	}, nil

}

func convertUpdatedUserToDB(pbEvent *pbV1Updated.Event) (*domain.User, error) {
	uuidDB, err := uuid.Parse(pbEvent.GetData().GetPublicId())
	if err != nil {
		return nil, fmt.Errorf("rpc user must have id uuid type: %w", err)
	}

	userRoleDB := convertRPCUserRoleToDB(pbEvent.GetData().GetRole())
	if userRoleDB == domain.UnknownRole {
		logrus.Errorf("rpc created user event: user have unknown type user role")

		return nil, fmt.Errorf("rpc created user event: user have unknown type user role")
	}

	return &domain.User{
		ID:        uuidDB,
		FirstName: pbEvent.GetData().GetFirstName(),
		LastName:  pbEvent.GetData().GetLastName(),
		Username:  pbEvent.GetData().GetUsername(),
		Email:     pbEvent.GetData().GetEmail(),
		Role:      userRoleDB,
	}, nil

}

func convertRPCUserRoleToDB(rpcRole pbV1UserRole.UserRole) domain.UserRoleType {
	switch rpcRole {
	case pbV1UserRole.UserRole_USER_ROLE_ADMIN:
		return domain.AdministratorRole
	case pbV1UserRole.UserRole_USER_ROLE_BOOKKEEPER:
		return domain.BookkeepperRole
	case pbV1UserRole.UserRole_USER_ROLE_MANAGER:
		return domain.ManagerRole
	case pbV1UserRole.UserRole_USER_ROLE_WORKER:
		return domain.WorkerRole
	case pbV1UserRole.UserRole_USER_ROLE_UNSPECIFIED:
		fallthrough
	default:
		return domain.UnknownRole
	}
}
