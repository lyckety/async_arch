package tasks

import (
	"context"
	"fmt"
	"log"
	"time"

	mbProd "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/mb/producer"
	pbV1Events "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/grpc/taskevents/v1"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type TaskCUDEventSender struct {
	client *mbProd.MBProducer

	topic     string
	partition int
}

func New(brokerURLs []string, topic string, partition int) *TaskCUDEventSender {
	sender := &TaskCUDEventSender{
		client:    mbProd.New(brokerURLs),
		topic:     topic,
		partition: partition,
	}

	return sender
}

func (s *TaskCUDEventSender) Send(
	ctx context.Context,
	events ...*pbV1Events.TaskEvent,
) error {
	msgs := make([]kafka.Message, len(events))

	for i, event := range events {
		binData, err := proto.Marshal(event)
		if err != nil {
			log.Fatalf("Failed to serialize message: %v", err)
		}
		eventMsg := kafka.Message{
			Topic:     s.topic,
			Partition: s.partition,
			Key:       []byte(event.GetTask().GetId()),
			Value:     binData,
			Time:      time.Unix(event.GetTimestamp(), 0),
		}

		msgs[i] = eventMsg
	}

	if err := s.client.SendMessages(context.Background(), msgs...); err != nil {
		return fmt.Errorf("kafkaClient.SendMessages(...): %s", err.Error())
	}

	logrus.Debugf(
		"Successfully sent task cud event message(s) (batch size %d): Partition: %d, Topic: %s",
		len(events),
		s.partition,
		s.topic,
	)

	return nil
}
