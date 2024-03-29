package producer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

const (
	DefaultConnectTimeout = 5 * time.Second
	DefaultReadTimeout    = 5 * time.Second
	DefaultWriteTimeout   = 5 * time.Second
	DefaultMaxAttempts    = 3
	DefaultNumPartition   = 1
	DefaultPartition      = 0
)

type MBProducer struct {
	Client  *kafka.Writer
	Brokers []string
}

func New(brokerURLs []string) *MBProducer {
	client := &kafka.Writer{
		Addr:         kafka.TCP(brokerURLs...),
		MaxAttempts:  DefaultMaxAttempts,
		WriteTimeout: DefaultWriteTimeout,
		ReadTimeout:  DefaultReadTimeout,
		Balancer:     &kafka.LeastBytes{},

		// Указывает гарантию доставки сообщения:
		// RequireAll - подтверждение лидера и реплики,
		// RequireOne - подтверждение только лидера
		// RequireNone - не подтверждать запись
		RequiredAcks: kafka.RequireAll,

		// Разрешает автосоздание топика
		AllowAutoTopicCreation: true,
	}

	mb := &MBProducer{
		Client: client,
	}

	return mb
}

func (p *MBProducer) checkLeaderExist(topicName string, numChecks int) error {
	for i := 0; i < numChecks; i++ {
		if i != 0 {
			time.Sleep(5 * time.Second)
		}

		if _, err := kafka.DialLeader(
			context.Background(), "tcp",
			p.Brokers[0],
			topicName,
			DefaultPartition,
		); err != nil {
			continue
		}

		logrus.Printf("Leader found for topic %s\n", topicName)

		return nil
	}

	return fmt.Errorf("could not find leader for topic %q", topicName)
}

func (p *MBProducer) SendMessage(ctx context.Context, msg kafka.Message) error {
	if err := p.Client.WriteMessages(ctx, msg); err != nil {
		if errors.Is(err, kafka.LeaderNotAvailable) || strings.Contains(err.Error(), "Leader Not Available") {
			if err := p.checkLeaderExist(msg.Topic, 3); err != nil {
				return fmt.Errorf("p.checkLeaderExist(...): %w", err)
			}

			if err := p.Client.WriteMessages(ctx, msg); err != nil {
				logrus.Errorf(
					"Producer failed send message: Key: %q, Value: %q, Topic: %q, Date: %s",
					msg.Key,
					msg.Value,
					msg.Topic,
					msg.Time,
				)
			} else {
				logrus.Infof(
					"Producer success send message: Key: %q, Value: %q, Topic: %q, Date: %s",
					msg.Key,
					msg.Value,
					msg.Topic,
					msg.Time,
				)

				return nil
			}
		}

		logrus.Errorf(
			"Producer failed send message: Key: %q, Value: %q, Topic: %q, Date: %s",
			msg.Key,
			msg.Value,
			msg.Topic,
			msg.Time,
		)

		return fmt.Errorf("p.producer.WriteMessages(...): %w", err)
	}

	logrus.Debugf(
		"Producer success send message: Key: %q, Value: %q, Topic: %q, Date: %s",
		msg.Key,
		msg.Value,
		msg.Topic,
		msg.Time,
	)

	return nil
}
