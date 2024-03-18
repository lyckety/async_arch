package consumer

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	pbV1Created "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/created/v1"
	pbV1UserCreated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/created/v1"
	pbV1Updated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/updated/v1"
	pbV1UserUpdated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/updated/v1"
	pbV1UserRole "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/userrole/v1"
	"github.com/lyckety/async_arch/popug_jira/services/accounting/internal/db/domain"
	mbCons "github.com/lyckety/async_arch/popug_jira/services/accounting/internal/mb/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type UsersCUDProcessor struct {
	mbConsClient *mbCons.MBConsumer

	storage domain.Repository
}

func NewUsersCUDProcessor(
	repository domain.Repository,
	brokerURLs []string,
	topic string,
	groupIDPrefix string,
	partition int,
) *UsersCUDProcessor {
	return &UsersCUDProcessor{
		mbConsClient: mbCons.New(brokerURLs, topic, fmt.Sprintf("%s-%s", groupIDPrefix, topic), partition),
		storage:      repository,
	}
}

func (p *UsersCUDProcessor) StartProcessMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		message, err := p.mbConsClient.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, mbCons.ErrContextCanceled) {
				logrus.Debug("p.mbConsClient.FetchMessage(ctx): context canceled")
			} else {
				logrus.Errorf("UserCUDEventProcessor: error fetch message from mb: %s", err)

				if err := p.mbConsClient.CommitMessages(ctx, message); err != nil {
					logrus.Errorf(
						"p.mbConsClient.CommitMessages(...): error commit message %v from mb: %s", message, err.Error(),
					)
				}
			}

			continue
		}

		logrus.Debugf("EventConsumeProcessor: success fetched message %v from mb", message)

		if err := p.processMessage(ctx, message); err != nil {
			logrus.Errorf(
				"p.processMessage(...): error commit message %v from mb: %s", message, err.Error(),
			)
		} else {
			logrus.Debugf("p.processMessage(...): success processed message %v from mb!", message)
		}

		if err := p.mbConsClient.CommitMessages(ctx, message); err != nil {
			logrus.Errorf("p.mbConsClient.CommitMessages(...): error commit message %v from mb: %s", message, err.Error())
		} else {
			logrus.Debugf("p.mbConsClient.CommitMessages(...): success commit message %v from mb", message)
		}

		logrus.Debugf("UserCUDEventProcessor: success processed message %v from mb!", message)
	}
}

func (p *UsersCUDProcessor) processMessage(ctx context.Context, msg kafka.Message) error {
	switch pbData := getUnmarshalledPbEventFromBinary(msg.Value).(type) {
	case *pbV1UserCreated.Event:
		dbUser, err := convertCreatedUserToDB(pbData)
		if err != nil {
			return fmt.Errorf("processMessage(msg): %v :%w", pbData, err)
		}

		if _, err := p.storage.CreateOrUpdateUser(ctx, dbUser); err != nil {
			return fmt.Errorf("processMessage(msg): error p.storage.CreateOrUpdateUser(ctx, %v): %w", dbUser, err)
		}

		return nil
	case *pbV1UserUpdated.Event:
		dbUser, err := convertUpdatedUserToDB(pbData)
		if err != nil {
			return fmt.Errorf("processMessage(msg): %q :%w", pbData, err)
		}

		if _, err := p.storage.CreateOrUpdateUser(ctx, dbUser); err != nil {
			return fmt.Errorf("processMessage(msg): error p.storage.CreateOrUpdateUser(ctx, %v): %w", dbUser, err)
		}
		return nil
	default:
		return fmt.Errorf("processMessage(msg): unknown protobuf type for event message")
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
		PublicID:  uuidDB,
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
		PublicID:  uuidDB,
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
