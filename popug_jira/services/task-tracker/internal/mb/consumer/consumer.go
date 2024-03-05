package consumer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

const (
	contextCanceled = "context canceled"
)

var ErrContextCanceled = errors.New(contextCanceled)

type Message struct {
	Key   []byte
	Value []byte
}

type MBConsumer struct {
	client *kafka.Reader

	brokers   []string
	partition int
	topic     string
	groupID   string
}

func New(
	brokers []string,
	topic string,
	groupID string,
	partition int,
) *MBConsumer {
	consumerCfg := kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,

		MinBytes: 0,    // 1B
		MaxBytes: 10e6, // 10MB

		Partition:   partition,
		StartOffset: kafka.FirstOffset,
		GroupID:     groupID,
	}

	return &MBConsumer{
		client:    kafka.NewReader(consumerCfg),
		brokers:   brokers,
		partition: partition,
		topic:     topic,
	}
}

func (c *MBConsumer) Close() {
	if err := c.client.Close(); err != nil {
		logrus.Errorf(
			"error close stream from mb consumer %s for topic %q and partition %q: %s",
			strings.Join(c.brokers, ", "),
			c.topic,
			c.partition,
			err.Error(),
		)
	}
}

func (c *MBConsumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	message, err := c.client.FetchMessage(ctx)
	if err != nil {
		if err.Error() == contextCanceled {
			logrus.Debug("p.producer.FetchMessage(...): context canceled")

			return kafka.Message{}, ErrContextCanceled
		}

		logrus.Errorf(
			"error fetch msg from mb consumer %s for topic %q and partition %q: %s",
			strings.Join(c.brokers, ", "),
			c.topic,
			c.partition,
			err.Error(),
		)

		return kafka.Message{}, fmt.Errorf("c.consumer.FetchMessage(...): %w", err)
	}

	logrus.Debugf(
		"success fetch msg %v from mb consumer %s for topic %q and partition %q!",
		message,
		strings.Join(c.brokers, ", "),
		c.topic,
		c.partition,
	)

	return message, nil
}

func (c *MBConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	message, err := c.client.ReadMessage(ctx)
	if err != nil {
		if err.Error() == contextCanceled {
			logrus.Debug("p.producer.ReadMessage(...): context canceled")

			return kafka.Message{}, ErrContextCanceled
		}

		logrus.Errorf(
			"error read msg from mb consumer %s for topic %q and partition %q: %s",
			strings.Join(c.brokers, ", "),
			c.topic,
			c.partition,
			err.Error(),
		)

		return kafka.Message{}, fmt.Errorf("c.consumer.ReadMessage(...): %w", err)
	}

	logrus.Debugf(
		"success read msg %v from mb consumer %s for topic %q and partition %q!",
		message,
		strings.Join(c.brokers, ", "),
		c.topic,
		c.partition,
	)

	return message, nil
}

func (c *MBConsumer) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	if err := c.client.CommitMessages(ctx, msgs...); err != nil {
		logrus.Errorf(
			"error commit msgs mb consumer %s for topic %q and partition %q: %s",
			strings.Join(c.brokers, ", "),
			c.topic,
			c.partition,
			err.Error(),
		)

		return fmt.Errorf("c.consumer.FetchMessage(...): %w", err)
	}

	logrus.Debugf(
		"success commit msgs mb consumer %s for topic %q and partition %q!",
		strings.Join(c.brokers, ", "),
		c.topic,
		c.partition,
	)

	return nil
}
