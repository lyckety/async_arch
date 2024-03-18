package consumer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	pbV1TaskAssigned "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/assigned/v1"
	pbV1TaskCompleted "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/completed/v1"
	prodEvents "github.com/lyckety/async_arch/popug_jira/services/accounting/internal/app/events/producer"
	"github.com/lyckety/async_arch/popug_jira/services/accounting/internal/db/domain"
	mbCons "github.com/lyckety/async_arch/popug_jira/services/accounting/internal/mb/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type TasksBEProcessor struct {
	mbConsClient *mbCons.MBConsumer
	mbProdClient *prodEvents.TxEventSender

	storage domain.Repository
}

func NewTasksBEProcessor(
	repository domain.Repository,
	brokerURLs []string,
	topic string,
	groupIDPrefix string,
	partition int,
	txProducer *prodEvents.TxEventSender,
) *TasksBEProcessor {
	return &TasksBEProcessor{
		mbConsClient: mbCons.New(brokerURLs, topic, fmt.Sprintf("%s-%s", groupIDPrefix, topic), partition),
		mbProdClient: txProducer,
		storage:      repository,
	}
}

func (p *TasksBEProcessor) StartProcessMessages(ctx context.Context) {
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

		logrus.Debugf("TasksBEProcessor: success processed message %v from mb!", message)
	}
}

func (p *TasksBEProcessor) processMessage(ctx context.Context, msg kafka.Message) error {
	switch pbData := getUnmarshalledPbEventFromBinary(msg.Value).(type) {
	case *pbV1TaskAssigned.Event:
		uuidTask, err := uuid.Parse(pbData.GetData().GetPublicId())
		if err != nil {
			return fmt.Errorf("rpc task must have id uuid type: %w", err)
		}

		uuidWorker, err := uuid.Parse(pbData.GetData().GetWorkerPublicId())
		if err != nil {
			return fmt.Errorf("rpc worker id must have id uuid type: %w", err)
		}

		taskDB, err := p.storage.GetTaskByPublicID(ctx, uuidTask)
		if err != nil {
			return fmt.Errorf("p.storage.GetTaskByPublicID(ctx, %s): %w", uuidTask.String(), err)
		}

		// Гонка между assigned и created решена как неполное создание задачи. Если created прилетел позже,
		// то обновит поле description, а если наоборот, то не обновится (поле не станет вдруг пустым)
		if taskDB == nil {
			taskDB = &domain.Task{
				PublicID: uuidTask,
			}

			taskDB.CostAssign, taskDB.CostComplete = generateTaskCosts()

			// if task already exist with helpful cost values not be updated
			_, err = p.storage.CreateOrUpdateTask(ctx, taskDB)
			if err != nil {
				return fmt.Errorf("error create task in db with id: %w", err)
			}
		}
		eventTime := time.Unix(pbData.GetHeader().GetEventTime(), 0)

		// TODO: retry if failed
		txDB, err := p.storage.CreateCreditTransaction(ctx, uuidTask, uuidWorker, eventTime)
		if err != nil {
			return fmt.Errorf("error create credit transaction: %w", err)
		}

		logrus.Debugf("created transaction credit %v", txDB)

		eventsMsg := prodEvents.NewTxCreditedV1(
			txDB.PublicID,
			uuidWorker,
			taskDB.PublicID,
			uint32(txDB.Credit),
			txDB.UpdatedAt.Unix(),
		)

		if err := p.mbProdClient.Send(context.Background(), eventsMsg); err != nil {
			logrus.Errorf("failed send cud event for tx credit: %s", err.Error())

			return fmt.Errorf("failed send cud event for tx credit: %w", err)
		}

		logrus.Debug("success sent cud event (tx credit)!")

		return nil
	case *pbV1TaskCompleted.Event:
		uuidTask, err := uuid.Parse(pbData.GetData().GetPublicId())
		if err != nil {
			return fmt.Errorf("rpc task must have id uuid type: %w", err)
		}

		uuidWorker, err := uuid.Parse(pbData.GetData().GetWorkerPublicId())
		if err != nil {
			return fmt.Errorf("rpc worker id must have id uuid type: %w", err)
		}

		taskDB, err := p.storage.GetTaskByPublicID(ctx, uuidTask)
		if err != nil {
			return fmt.Errorf("p.storage.GetTaskByPublicID(ctx, %s): %w", uuidTask.String(), err)
		}

		// Гонка между completed и created решена как неполное создание задачи. Если created прилетел позже,
		// то обновит поле description, а если наоборот, то не обновится (поле не станет вдруг пустым)
		if taskDB == nil {
			taskDB = &domain.Task{
				PublicID: uuidTask,
			}

			taskDB.CostAssign, taskDB.CostComplete = generateTaskCosts()

			// if task already exist with helpful cost values not be updated
			_, err = p.storage.CreateOrUpdateTask(ctx, taskDB)
			if err != nil {
				return fmt.Errorf("error create task in db with id: %w", err)
			}
		}

		eventTime := time.Unix(pbData.GetHeader().GetEventTime(), 0)

		// TODO: retry if failed
		txDB, err := p.storage.CreateDebitTransaction(ctx, uuidTask, uuidWorker, eventTime)
		if err != nil {
			return fmt.Errorf("error create debit transaction: %w", err)
		}

		logrus.Debugf("created transaction debit %v", txDB)

		eventsMsg := prodEvents.NewTxDebitedV1(
			txDB.PublicID,
			txDB.UserPublicID,
			txDB.TaskPublicID,
			uint32(txDB.Debit),
			txDB.UpdatedAt.Unix(),
		)

		if err := p.mbProdClient.Send(context.Background(), eventsMsg); err != nil {
			logrus.Errorf("failed send cud event for tx debit: %s", err.Error())

			return fmt.Errorf("failed send cud event for tx dbit: %w", err)
		}

		logrus.Debug("success sent cud event (tx debit)!")

		return nil
	default:
		return fmt.Errorf("processMessage(msg): unknown protobuf type for event message")
	}
}
