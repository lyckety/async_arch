package consumer

import (
	"context"
	"fmt"

	"github.com/lyckety/async_arch/golang-service/internal/kafka/consumer/handler"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type HandlerFetchMessage func(message kafka.Message) error

type Consumer struct {
	brokers            []string
	topic              string
	config             *kafka.ReaderConfig
	consumer           *kafka.Reader
	consumerType       string
	weightFailIfFailed float32
}

func New(
	brokers []string,
	topic,
	groupID string,
	partition int,
	readerType string,
) (*Consumer, error) {
	if groupID != "" && partition > 0 {
		return nil, fmt.Errorf("not allowed set partition and groupID for consumer")
	}

	p := &Consumer{
		brokers:      brokers,
		topic:        topic,
		consumerType: readerType,
	}

	p.config = &kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,

		MinBytes: 0,    // 1B
		MaxBytes: 10e6, // 10MB

		// можно указывать ТОЛЬКО или партицию, или groupID
		Partition:   partition,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	}

	p.consumer = kafka.NewReader(*p.config)

	return p, nil
}

func (c *Consumer) Close() {
	c.consumer.Close()
}

func (c *Consumer) Consume(ctx context.Context) error {
	switch c.consumerType {
	case "fetch":
		if err := c.fetchMessages(ctx); err != nil {
			return fmt.Errorf("c.fetchMessages(...): %q", err)
		}

		return nil
	case "read":
		if err := c.readMessages(ctx); err != nil {
			return fmt.Errorf("c.readMessages(...): %q", err)
		}

		return nil
	default:
		return fmt.Errorf("not correct consume type : %q", c.consumerType)
	}
}

func (c *Consumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	message, err := c.consumer.FetchMessage(ctx)
	if err != nil {
		if err.Error() == "context canceled" {
			logrus.Debug("c.consumer.ReadMessage(...): context canceled")

			return kafka.Message{}, nil
		}

		logrus.Errorf("Ошибка чтения сообщения : %v\n", err)

		return kafka.Message{}, fmt.Errorf("c.consumer.ReadMessage(...): %w", err)
	}

	return message, nil
}

// Данный метод делает автоматический commit offset, не подойдет если у нас есть обработчик, который может сломаться и придется перечитать сообщение
func (c *Consumer) readMessages(ctx context.Context) error {
	for {
		message, err := c.consumer.ReadMessage(ctx)
		if err != nil {
			if err.Error() == "context canceled" {
				logrus.Debug("c.consumer.ReadMessage(...): context canceled")

				return nil
			}

			logrus.Errorf("Ошибка чтения сообщения : %v\n", err)

			return fmt.Errorf("c.consumer.ReadMessage(...): %w", err)
		}

		logrus.Infof(
			"Consumer успешно получил сообщение : Key: %q, Value: %q, Topic: %q, Date: %s",
			message.Key,
			message.Value,
			c.consumer.Config().Topic,
			message.Time,
		)

		if err := handler.Handle(message, c.weightFailIfFailed); err != nil {
			logrus.Errorf(
				"Ошибка обработки сообщения : Key: %q, Value: %q, Topic: %q, Date: %s",
				message.Key,
				message.Value,
				c.consumer.Config().Topic,
				message.Time,
			)

		}

		logrus.Infof(
			"Сообщение успешно получено и обработано: Key: %q, Value: %q, Topic: %q, Date: %s",
			message.Key,
			message.Value,
			c.consumer.Config().Topic,
			message.Time,
		)
	}
}

func (c *Consumer) fetchMessages(ctx context.Context) error {
	for {
		message, err := c.consumer.FetchMessage(ctx)
		if err != nil {
			if err.Error() == "context canceled" {
				logrus.Debug("c.consumer.ReadMessage(...): context canceled")

				return nil
			}

			logrus.Errorf("Ошибка чтения сообщения : %v\n", err)

			return fmt.Errorf("c.consumer.ReadMessage(...): %w", err)
		}

		logrus.Infof(
			"Consumer успешно получил сообщение : Key: %q, Value: %q, Topic: %q, Date: %s",
			message.Key,
			message.Value,
			c.consumer.Config().Topic,
			message.Time,
		)

		if err := handler.Handle(message, c.weightFailIfFailed); err != nil {
			logrus.Errorf(
				"Ошибка обработки сообщения : Key: %q, Value: %q, Topic: %q, Date: %s",
				message.Key,
				message.Value,
				c.consumer.Config().Topic,
				message.Time,
			)

			return fmt.Errorf("c.consumer.ReadMessage(...): %w", err)

		}

		if err := c.consumer.CommitMessages(ctx, message); err != nil {
			logrus.Errorf(
				"Ошибка обработки сообщения : Key: %q, Value: %q, Topic: %q, Date: %s",
				message.Key,
				message.Value,
				c.consumer.Config().Topic,
				message.Time,
			)

			return fmt.Errorf("c.consumer.ReadMessage(...): %w", err)

		}

		logrus.Infof(
			"Сообщение успешно получено и обработано: Key: %q, Value: %q, Topic: %q, Date: %s",
			message.Key,
			message.Value,
			c.consumer.Config().Topic,
			message.Time,
		)
	}
}

func (c *Consumer) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	if err := c.consumer.CommitMessages(ctx, msgs...); err != nil {
		if err.Error() == "context canceled" {
			logrus.Debug("c.consumer.CommitMessages(...): context canceled")

			return nil
		}

		return fmt.Errorf("c.consumer.CommitMessages(...): %w", err)
	}

	logrus.Info("successfull commited messages!")

	return nil
}
