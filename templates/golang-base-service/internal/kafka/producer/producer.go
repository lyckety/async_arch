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
	DefaultConnectTimeout = 3 * time.Second
	DefaultReadTimeout    = 3 * time.Second
	DefaultWriteTimeout   = 3 * time.Second
	DefaultMaxAttempts    = 3
	DefaultNumPartition   = 1
	DefaultPartition      = 0
)

type Producer struct {
	brokers []string

	producer *kafka.Writer
}

func New(brokers []string) *Producer {
	p := &Producer{
		brokers: brokers,
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(p.brokers...),
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

	p.producer = w

	return p
}

func (p *Producer) Close() {
	p.producer.Close()
}

func (p *Producer) SendMessages(ctx context.Context, messages ...kafka.Message) error {
	for _, msg := range messages {
		if err := p.producer.WriteMessages(ctx, msg); err != nil {
			if err.Error() == "context canceled" {
				logrus.Debug("p.producer.WriteMessages(...): context canceled")

				return nil
			}

			if errors.Is(err, kafka.LeaderNotAvailable) || strings.Contains(err.Error(), "Leader Not Available") {
				if err := p.checkLeaderExist(msg.Topic, 3); err != nil {
					return fmt.Errorf("p.checkLeaderExist(...): %w", err)
				}

				if err := p.producer.WriteMessages(ctx, msg); err != nil {
					if err.Error() == "context canceled" {
						logrus.Debug("p.producer.WriteMessages(...): context canceled")

						return nil
					}

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

		logrus.Infof(
			"Producer success send message: Key: %q, Value: %q, Topic: %q, Date: %s",
			msg.Key,
			msg.Value,
			msg.Topic,
			msg.Time,
		)
	}

	return nil
}

func (p *Producer) checkLeaderExist(topicName string, numChecks int) error {
	for i := 0; i < numChecks; i++ {
		if i != 0 {
			time.Sleep(5 * time.Second)
		}

		if _, err := kafka.DialLeader(
			context.Background(), "tcp",
			p.brokers[0],
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
