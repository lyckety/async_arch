package consumer

import (
	"context"
	"errors"
	"fmt"

	pbV1TaskCreated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/created/v1"
	mbCons "github.com/lyckety/async_arch/popug_jira/services/analytics/internal/mb/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type TasksStreamProcessor struct {
	mbConsClient *mbCons.MBConsumer
}

func NewTasksStreamProcessor(
	brokerURLs []string,
	topic string,
	groupIDPrefix string,
	partition int,
) *TasksStreamProcessor {
	return &TasksStreamProcessor{
		mbConsClient: mbCons.New(brokerURLs, topic, fmt.Sprintf("%s-%s", groupIDPrefix, topic), partition),
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
				logrus.Errorf("TasksStreamProcessor: error fetch message from mb: %s", err)

				if err := p.mbConsClient.CommitMessages(ctx, message); err != nil {
					logrus.Errorf(
						"p.mbConsClient.CommitMessages(...): error commit message from mb: %s", err.Error(),
					)
				}
			}

			continue
		}

		logrus.Debugf("EventConsumeProcessor: success fetched message from mb")

		if err := p.processMessage(ctx, message); err != nil {
			logrus.Errorf(
				"p.processMessage(...): error commit message from mb: %s", err.Error(),
			)
		} else {
			logrus.Debug("p.processMessage(...): success processed message from mb!")
		}

		if err := p.mbConsClient.CommitMessages(ctx, message); err != nil {
			logrus.Errorf("p.mbConsClient.CommitMessages(...): error commit message from mb: %s", err.Error())
		} else {
			logrus.Debug("p.mbConsClient.CommitMessages(...): success commit message from mb")
		}

		logrus.Debug("TasksStreamProcessor: success processed message from mb!")
	}
}

func (p *TasksStreamProcessor) processMessage(ctx context.Context, msg kafka.Message) error {
	switch pbData := getUnmarshalledPbEventFromBinary(msg.Value).(type) {
	case *pbV1TaskCreated.Event:
		logrus.Infof("============== Fetched TASKS stream event (credited): %v ==============", pbData)

		return nil
	default:
		return fmt.Errorf("processMessage(msg): unknown protobuf type for event message")
	}
}
