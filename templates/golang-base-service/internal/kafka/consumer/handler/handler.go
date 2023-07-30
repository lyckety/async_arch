package handler

import (
	"fmt"
	"math/rand"
	"time"

	wr "github.com/mroth/weightedrand"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

const (
	defaultMultipleForWeight = 10
)

func Handle(msg kafka.Message, probabilityErr float32) error {
	rand.Seed(time.Now().Unix())

	for _, header := range msg.Headers {
		if header.Key == "isFailed" {
			logrus.Infof(
				"handler: marked as failed message: Key: %q, Value: %q, Topic: %q, Date: %v, calculate probability error...",
				msg.Key,
				msg.Value,
				msg.Topic,
				msg.Time,
			)

			isFailed := isItPossibleByWeight(probabilityErr)

			if isFailed {
				logrus.Errorf(
					"handler: failed handle message: Key: %q, Value: %q, Topic: %q, Date: %v",
					msg.Key,
					msg.Value,
					msg.Topic,
					msg.Time,
				)

				return fmt.Errorf(
					"handler: failed handle message: Key: %q, Value: %q, Topic: %q, Date: %v",
					msg.Key,
					msg.Value,
					msg.Topic,
					msg.Time,
				)
			}
		}
	}

	logrus.Infof(
		"handler: succes handle message: Key: %q, Value: %q, Topic: %q, Date: %v",
		msg.Key,
		msg.Value,
		msg.Topic,
		msg.Time,
	)

	return nil
}

func isItPossibleByWeight(weight float32) bool {
	rand.Seed(time.Now().Unix())

	possibility := uint(defaultMultipleForWeight * weight)

	notPossibility := defaultMultipleForWeight - possibility

	chooser, _ := wr.NewChooser(
		wr.Choice{Item: false, Weight: notPossibility},
		wr.Choice{Item: true, Weight: possibility},
	)

	return chooser.Pick().(bool) //nolint: forcetypeassert
}
