package events

import (
	"context"
	"fmt"
	"log"
	"time"

	mbProd "github.com/lyckety/async_arch/popug_jira/services/auth/internal/mb/producer"
	pbV1Events "github.com/lyckety/async_arch/popug_jira/services/auth/pkg/grpc/authevents/v1"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type CUDEventSender struct {
	client *mbProd.MBProducer

	topic     string
	partition int
}

func New(brokerURLs []string, topic string, partition int) *CUDEventSender {
	sender := &CUDEventSender{
		client:    mbProd.New(brokerURLs),
		topic:     topic,
		partition: partition,
	}

	return sender
}

func (s *CUDEventSender) Send(
	ctx context.Context,
	event *pbV1Events.AuthEvent,
) error {
	binData, err := proto.Marshal(event)
	if err != nil {
		log.Fatalf("Failed to serialize message: %v", err)
	}
	eventMsg := kafka.Message{
		Topic:     s.topic,
		Partition: s.partition,
		Key:       []byte(event.GetUser().GetId()),
		Value:     binData,
		Time:      time.Unix(event.GetTimestamp(), 0),
	}

	if err := s.client.SendMessage(ctx, eventMsg); err != nil {
		logrus.Errorf("kafkaClient.SendMessages(...): %s", err.Error())

		return fmt.Errorf("kafkaClient.SendMessages(...): %s", err.Error())
	}

	logrus.Debugf(
		"Successfully sent cud event message: Key: %s, Value: %s, Partition: %d, Topic: %s, Time: %s",
		fmt.Sprint(eventMsg.Key),
		eventMsg.Value,
		s.partition,
		s.topic,
		eventMsg.Time,
	)

	return nil
}
