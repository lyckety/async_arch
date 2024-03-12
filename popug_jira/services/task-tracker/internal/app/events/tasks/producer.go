package tasks

import (
	"context"
	"errors"
	"fmt"
	"time"

	mbProd "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/mb/producer"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrToType = errors.New("failed convert to type")
)

type TaskEventSender struct {
	client *mbProd.MBProducer

	topic     string
	partition int
}

func New(brokerURLs []string, topic string, partition int) *TaskEventSender {
	sender := &TaskEventSender{
		client:    mbProd.New(brokerURLs),
		topic:     topic,
		partition: partition,
	}

	return sender
}

func (s *TaskEventSender) Send(
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
		case *TaskCreatedV1:
			logrus.Debug("Start processing write event TaskCreatedV1")
			e, ok := event.(*TaskCreatedV1)
			if !ok {
				return fmt.Errorf("%w TaskCreatedV1", ErrToType)
			}

			pbMsg, err = e.ToPB()
			if err != nil {
				return fmt.Errorf("error convert TaskCreatedV1 to protobuf event: %w", err)
			}

			kafkaMsg.Key = []byte(e.PublicID.String())
			kafkaMsg.Time = time.Unix(e.EventTime, 0)
		case *TaskAssignedV1:
			logrus.Debug("Start processing write event TaskAssignedV1")
			e, ok := event.(*TaskAssignedV1)
			if !ok {
				return fmt.Errorf("%w TaskAssignedV1", ErrToType)
			}

			pbMsg, err = e.ToPB()
			if err != nil {
				return fmt.Errorf("error convert TaskAssignedV1 to protobuf event: %w", err)
			}

			kafkaMsg.Key = []byte(e.PublicID.String())
			kafkaMsg.Time = time.Unix(e.EventTime, 0)
		case *TaskCompletedV1:
			logrus.Debug("Start processing write event TaskCompletedV1")
			e, ok := event.(*TaskCompletedV1)
			if !ok {
				return fmt.Errorf("%w TaskCompletedV1", ErrToType)
			}

			pbMsg, err = e.ToPB()
			if err != nil {
				return fmt.Errorf("error convert TaskCompletedV1 to protobuf event: %w", err)
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
