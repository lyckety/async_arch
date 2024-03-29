package tasks

import (
	"context"
	"errors"
	"fmt"
	"time"

	mbProd "github.com/lyckety/async_arch/popug_jira/services/accounting/internal/mb/producer"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrToType = errors.New("failed convert to type")
)

type TxEventSender struct {
	client *mbProd.MBProducer

	topic     string
	partition int
}

func New(brokerURLs []string, topic string, partition int) *TxEventSender {
	sender := &TxEventSender{
		client:    mbProd.New(brokerURLs),
		topic:     topic,
		partition: partition,
	}

	return sender
}

func (s *TxEventSender) Send(
	ctx context.Context,
	events ...interface{},
) error {
	var err error

	kafkaMsgs := make([]kafka.Message, 0)

	for _, event := range events {
		var pbMsg protoreflect.ProtoMessage

		var kafkaMsg = kafka.Message{
			Topic: s.topic,
		}

		switch et := event.(type) {
		case *TxDebitedV1:
			logrus.Debug("Start processing write event TxDebitedV1")
			e, ok := event.(*TxDebitedV1)
			if !ok {
				return fmt.Errorf("%w TxDebitedV1", ErrToType)
			}

			pbMsg, err = e.ToPB()
			if err != nil {
				return fmt.Errorf("error convert TxDebitedV1 to protobuf event: %w", err)
			}

			kafkaMsg.Key = []byte(e.PublicID.String())
			kafkaMsg.Time = time.Unix(e.EventTime, 0)
		case *TxCreditedV1:
			logrus.Debug("Start processing write event TxCreditedV1")
			e, ok := event.(*TxCreditedV1)
			if !ok {
				return fmt.Errorf("%w TxCreditedV1", ErrToType)
			}

			pbMsg, err = e.ToPB()
			if err != nil {
				return fmt.Errorf("error convert TxCreditedV1 to protobuf event: %w", err)
			}

			kafkaMsg.Key = []byte(e.PublicID.String())
			kafkaMsg.Time = time.Unix(e.EventTime, 0)
		case *TxPaymentedV1:
			logrus.Debug("Start processing write event TxPaymentedV1")
			e, ok := event.(*TxPaymentedV1)
			if !ok {
				return fmt.Errorf("%w TxPaymentedV1", ErrToType)
			}

			pbMsg, err = e.ToPB()
			if err != nil {
				return fmt.Errorf("error convert TxPaymentedV1 to protobuf event: %w", err)
			}

			kafkaMsg.Key = []byte(e.PublicID.String())
			kafkaMsg.Time = time.Unix(e.EventTime, 0)
		default:
			logrus.Errorf("send task events: unknown type event: %v", et)

			return fmt.Errorf("send task events: unknown type event: %v", et)
		}

		kafkaMsg.Value, err = proto.Marshal(pbMsg)
		if err != nil {
			return fmt.Errorf("failed to serialize protobuf message: %v: %w", pbMsg, err)
		}

		kafkaMsgs = append(kafkaMsgs, kafkaMsg)
	}

	if err := s.client.SendMessages(ctx, kafkaMsgs...); err != nil {
		logrus.Errorf("kafkaClient.SendMessages(...): %s", err.Error())

		return fmt.Errorf("kafkaClient.SendMessages(...): %s", err.Error())
	}

	logrus.Debugf(
		"Successfully sent event messages (%d) from (%d) to topic: %s",
		len(kafkaMsgs),
		len(events),
		s.topic,
	)

	return nil
}
