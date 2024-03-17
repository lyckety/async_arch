package consumer

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	pbV1TaskCreated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/created/v1"
	pbV2TaskCreated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/created/v2"
	"github.com/lyckety/async_arch/popug_jira/services/accounting/internal/db/domain"
	mbCons "github.com/lyckety/async_arch/popug_jira/services/accounting/internal/mb/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type TasksStreamProcessor struct {
	mbConsClient *mbCons.MBConsumer

	storage domain.Repository
}

func NewTasksStreamProcessor(
	repository domain.Repository,
	brokerURLs []string,
	topic string,
	groupIDPrefix string,
	partition int,
) *TasksStreamProcessor {
	return &TasksStreamProcessor{
		mbConsClient: mbCons.New(brokerURLs, topic, fmt.Sprintf("%s-%s", groupIDPrefix, topic), partition),
		storage:      repository,
	}
}

func (p *TasksStreamProcessor) StartProcessMessages(ctx context.Context) {
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

func (p *TasksStreamProcessor) processMessage(ctx context.Context, msg kafka.Message) error {
	switch pbData := getUnmarshalledPbEventFromBinary(msg.Value).(type) {
	case *pbV1TaskCreated.Event:
		uuidDB, err := uuid.Parse(pbData.GetData().GetPublicId())
		if err != nil {
			return fmt.Errorf("rpc task must have id uuid type: %w", err)
		}

		jiraID, title := parseTaskDescription(pbData.GetData().GetDescription())

		dbTask := &domain.Task{
			PublicID:     uuidDB,
			JiraId:       jiraID,
			Description:  title,
			CostAssign:   0,
			CostComplete: 0,
		}

		dbTask.CostAssign, dbTask.CostComplete = generateTaskCosts()

		// if task already exist with helpful cost values not be updated
		_, err = p.storage.CreateOrUpdateTask(ctx, dbTask)
		if err != nil {
			return fmt.Errorf("error create task in db with id: %w", err)
		}

		logrus.Debugf("EventConsumeProcessor: success process event-message %q", pbData.GetHeader().GetEventName())

		return nil
	case *pbV2TaskCreated.Event:
		uuidDB, err := uuid.Parse(pbData.GetData().GetPublicId())
		if err != nil {
			return fmt.Errorf("rpc task must have id uuid type: %w", err)
		}

		// TODO: пока просто вывожу в консоль если  jira id и title некорректные
		if err := validateJiraIDAndTitle(pbData.GetData().GetJiraId(), pbData.GetData().GetDescription()); err != nil {
			logrus.Errorf("validate jira id and title in pbV2TaskCreated: %s", err.Error())
		}

		dbTask := &domain.Task{
			PublicID:     uuidDB,
			Description:  pbData.GetData().GetDescription(),
			JiraId:       pbData.GetData().GetJiraId(),
			CostAssign:   0,
			CostComplete: 0,
		}

		dbTask.CostAssign, dbTask.CostComplete = generateTaskCosts()

		// if task already exist with helpful cost values not be updated
		_, err = p.storage.CreateOrUpdateTask(ctx, dbTask)
		if err != nil {
			return fmt.Errorf("error create task in db with id: %w", err)
		}

		logrus.Debugf("EventConsumeProcessor: success process event-message %q", pbData.GetHeader().GetEventName())

		return nil
	default:
		return fmt.Errorf("processMessage(msg): unknown protobuf type for event message")
	}
}
