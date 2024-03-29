package consumer

import (
	"context"
	"errors"
	"fmt"

	pbV1TxCredited "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/transactionevents/credited/v1"
	pbV1TxDebited "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/transactionevents/debited/v1"
	pbV1TxPaymented "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/transactionevents/paymented/v1"
	mbCons "github.com/lyckety/async_arch/popug_jira/services/analytics/internal/mb/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type TxCUDProcessor struct {
	mbConsClient *mbCons.MBConsumer
}

func NewTxCUDProcessor(
	brokerURLs []string,
	topic string,
	groupIDPrefix string,
	partition int,
) *TxCUDProcessor {
	return &TxCUDProcessor{
		mbConsClient: mbCons.New(brokerURLs, topic, fmt.Sprintf("%s-%s", groupIDPrefix, topic), partition),
	}
}

func (p *TxCUDProcessor) StartProcessMessages(ctx context.Context) {
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
				logrus.Errorf("TxCUDProcessor: error fetch message from mb: %s", err)

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
			logrus.Debug("p.processMessage(...): success processed message  from mb!")
		}

		if err := p.mbConsClient.CommitMessages(ctx, message); err != nil {
			logrus.Errorf("p.mbConsClient.CommitMessages(...): error commit message from mb: %s", err.Error())
		} else {
			logrus.Debug("p.mbConsClient.CommitMessages(...): success commit message from mb")
		}

		logrus.Debug("TxCUDProcessor: success processed message from mb!")
	}
}

func (p *TxCUDProcessor) processMessage(ctx context.Context, msg kafka.Message) error {
	switch pbData := getUnmarshalledPbEventFromBinary(msg.Value).(type) {
	case *pbV1TxCredited.Event:
		logrus.Infof("============== Fetched TRANSACTIONS stream event (credited): %v ==============", pbData)

		return nil
	case *pbV1TxDebited.Event:
		logrus.Infof("============== Fetched TRANSACTIONS stream event (debited): %v ==============", pbData)

		return nil
	case *pbV1TxPaymented.Event:
		logrus.Infof("============== Fetched TRANSACTIONS stream event (paymented): %v ==============", pbData)

		return nil
	default:
		return fmt.Errorf("processMessage(msg): unknown protobuf type for event message %v", pbData)
	}
}
