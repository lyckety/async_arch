package consumer

import (
	"context"
	"errors"
	"fmt"

	pbV1Created "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/created/v1"
	pbV1Updated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/updated/v1"
	mbCons "github.com/lyckety/async_arch/popug_jira/services/analytics/internal/mb/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type UsersCUDProcessor struct {
	mbConsClient *mbCons.MBConsumer
}

func NewUsersCUDProcessor(
	brokerURLs []string,
	topic string,
	groupIDPrefix string,
	partition int,
) *UsersCUDProcessor {
	return &UsersCUDProcessor{
		mbConsClient: mbCons.New(brokerURLs, topic, fmt.Sprintf("%s-%s", groupIDPrefix, topic), partition),
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
						"p.mbConsClient.CommitMessages(...): error commit message from mb: %s", err.Error(),
					)
				}
			}

			continue
		}

		logrus.Debug("EventConsumeProcessor: success fetched message from mb")

		if err := p.processMessage(ctx, message); err != nil {
			logrus.Errorf(
				"p.processMessage(...): error commit message from mb: %s", err.Error(),
			)
		} else {
			logrus.Debug("p.processMessage(...): success processed message from mb!")
		}

		if err := p.mbConsClient.CommitMessages(ctx, message); err != nil {
			logrus.Errorf("p.mbConsClient.CommitMessages(...): error commit message %v from mb: %s", message, err.Error())
		} else {
			logrus.Debug("p.mbConsClient.CommitMessages(...): success commit message from mb")
		}

		logrus.Debug("UserCUDEventProcessor: success processed message from mb!")
	}
}

func (p *UsersCUDProcessor) processMessage(ctx context.Context, msg kafka.Message) error {
	switch pbData := getUnmarshalledPbEventFromBinary(msg.Value).(type) {
	case *pbV1Created.Event:
		logrus.Infof("============== Fetched USERS stream event (created): %v ==============", pbData)

		return nil
	case *pbV1Updated.Event:
		logrus.Infof("============== Fetched USERS stream event (updated): %v ==============", pbData)

		return nil
	default:
		return fmt.Errorf("processMessage(msg): unknown protobuf type for event message %v", pbData)
	}
}
