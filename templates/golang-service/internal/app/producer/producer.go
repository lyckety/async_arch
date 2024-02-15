package producer

import (
	"context"

	"time"

	kafkaProd "github.com/lyckety/async_arch/golang-service/internal/kafka/producer"
	pbV1 "github.com/lyckety/async_arch/golang-service/pkg/grpc/kafkaproducer/v1"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProducerService struct {
	pbV1.UnimplementedKafkaProducerServiceServer

	brokerURLs []string
}

func New(brokerURLs []string) *ProducerService {
	return &ProducerService{
		brokerURLs: brokerURLs,
	}
}

func (s *ProducerService) Send(
	ctx context.Context,
	req *pbV1.SendRequest,
) (*pbV1.SendResponse, error) {
	kafkaClient := kafkaProd.New(s.brokerURLs)
	defer kafkaClient.Close()

	kafkaMsg := rpcMsgToKafkaMsg(req)

	if err := kafkaClient.SendMessages(ctx, kafkaMsg); err != nil {
		return nil, status.Errorf(
			codes.Unknown,
			"kafkaClient.SendMessages(...): %s", err.Error(),
		)
	}

	logrus.Infof(
		"Successfully sent kafka message: Key: %s, Value: %s, Topic: %s, Time: %s, headers: %v",
		kafkaMsg.Key,
		kafkaMsg.Value,
		kafkaMsg.Topic,
		kafkaMsg.Time,
		kafkaMsg.Headers,
	)

	return &pbV1.SendResponse{}, nil
}

func rpcMsgToKafkaMsg(rpcMsg *pbV1.SendRequest) kafka.Message {
	tm := time.Unix(int64(rpcMsg.Timestamp), 0)

	kafkaMsg := kafka.Message{
		Topic: rpcMsg.GetTopic(),
		Key:   []byte(rpcMsg.GetKey()),
		Value: []byte(rpcMsg.GetValue()),
		Time:  tm,
	}

	return kafkaMsg
}
