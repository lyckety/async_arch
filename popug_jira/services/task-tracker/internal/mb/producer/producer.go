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

func (p *MBProducer) SendMessages(ctx context.Context, msgs ...kafka.Message) error {
	if err := p.Client.WriteMessages(ctx, msgs...); err != nil {
		if err.Error() == "context canceled" {
			logrus.Debug("p.producer.WriteMessages(...): context canceled")

			return nil
		}

		if errors.Is(err, kafka.LeaderNotAvailable) || strings.Contains(err.Error(), "Leader Not Available") {
			if err := p.checkLeaderExist(msgs[0].Topic, 3); err != nil {
				return fmt.Errorf("p.checkLeaderExist(...): %w", err)
			}

			if err := p.Client.WriteMessages(ctx, msgs[0]); err != nil {
				if err.Error() == "context canceled" {
					logrus.Debug("p.producer.WriteMessages(...): context canceled")

					return nil
				}

				logrus.Errorf(
					"Producer failed send message: Key: %q, Value: %q, Topic: %q, Date: %s",
					msgs[0].Key,
					msgs[0].Value,
					msgs[0].Topic,
					msgs[0].Time,
				)
			} else {
				logrus.Infof(
					"Producer success send message: Key: %q, Value: %q, Topic: %q, Date: %s",
					msgs[0].Key,
					msgs[0].Value,
					msgs[0].Topic,
					msgs[0].Time,
				)

				return nil
			}
		}

		logrus.Errorf(
			"Producer failed send message: Key: %q, Value: %q, Topic: %q, Date: %s",
			msgs[0].Key,
			msgs[0].Value,
			msgs[0].Topic,
			msgs[0].Time,
		)

		return fmt.Errorf("p.producer.WriteMessages(...): %w", err)
	}

	logrus.Debugf(
		"Producer success send message: Key: %q, Value: %q, Topic: %q, Date: %s",
		msgs[0].Key,
		msgs[0].Value,
		msgs[0].Topic,
		msgs[0].Time,
	)

	return nil
}
