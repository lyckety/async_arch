package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	mbProd "github.com/lyckety/async_arch/popug_jira/services/auth/internal/mb/producer"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrToType = errors.New("failed convert to type")
)

type UserCUDEventSender struct {
	client *mbProd.MBProducer

	topic     string
	partition int
}

func New(brokerURLs []string, topic string, partition int) *UserCUDEventSender {
	sender := &UserCUDEventSender{
		client:    mbProd.New(brokerURLs),
		topic:     topic,
		partition: partition,
	}

	return sender
}

func (s *UserCUDEventSender) Send(
	ctx context.Context,
	event interface{},
) error {
	var err error
	var kafkaMsg = kafka.Message{
		Topic: s.topic,
	}

	var pbMsg protoreflect.ProtoMessage

	switch et := event.(type) {
	case CreatedUserV1:
		logrus.Debug("Start processing write event CreatedUserV1")
		e, ok := event.(CreatedUserV1)
		if !ok {
			return fmt.Errorf("%w CreatedUserV1", ErrToType)
		}

		pbMsg, err = e.toPB()
		if err != nil {
			return fmt.Errorf("error convert CreatedUserV1 to protobuf event: %w", err)
		}

		kafkaMsg.Key = []byte(e.PublicID.String())
		kafkaMsg.Time = time.Unix(e.EventTime, 0)
	case UpdatedUserV1:
		logrus.Debug("Start processing write event UpdatedUserV1")
		e, ok := event.(UpdatedUserV1)
		if !ok {
			return fmt.Errorf("%w UpdatedUserV1", ErrToType)
		}

		pbMsg, err = e.toPB()
		if err != nil {
			return fmt.Errorf("error convert UpdatedUserV1 to protobuf event: %w", err)
		}

		kafkaMsg.Key = []byte(e.PublicID.String())
		kafkaMsg.Time = time.Unix(e.EventTime, 0)
	default:
		logrus.Errorf("send user cud events: unknown type event: %v", et)

		return fmt.Errorf("send user cud events: unknown type event: %v", et)
	}

	kafkaMsg.Value, err = proto.Marshal(pbMsg)
	if err != nil {
		return fmt.Errorf("failed to serialize protobuf message: %v: %w", pbMsg, err)
	}

	if err := s.client.SendMessage(ctx, kafkaMsg); err != nil {
		logrus.Errorf("kafkaClient.SendMessages(...): %s", err.Error())

		return fmt.Errorf("kafkaClient.SendMessages(...): %s", err.Error())
	}

	logrus.Debugf(
		"Successfully sent cud event message: Key: %s, Value: %s, Partition: %d, Topic: %s, Time: %s",
		fmt.Sprint(kafkaMsg.Key),
		kafkaMsg.Value,
		s.partition,
		s.topic,
		kafkaMsg.Time,
	)

	return nil
}
