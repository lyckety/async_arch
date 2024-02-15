package consumer

import (
	"context"

	kafkaCons "github.com/lyckety/async_arch/golang-service/internal/kafka/consumer"
	pbV1 "github.com/lyckety/async_arch/golang-service/pkg/grpc/kafkaconsumer/v1"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConsumerService struct {
	pbV1.UnimplementedKafkaConsumerServiceServer

	brokerURLs []string
}

func New(brokerURLs []string) *ConsumerService {
	return &ConsumerService{
		brokerURLs: brokerURLs,
	}
}

func (s *ConsumerService) SubscribeTopic(
	req *pbV1.SubscribeTopicRequest,
	stream pbV1.KafkaConsumerService_SubscribeTopicServer,
) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	kafkaClient, err := kafkaCons.New(s.brokerURLs, req.GetTopic(), "", 0, "fetch")
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"kafkaCons.New(s.brokerURLs, req.GetTopic(), \"\", 1, \"fetch\"): %s", err.Error(),
		)
	}
	defer kafkaClient.Close()

	for {
		msg, err := kafkaClient.FetchMessage(ctx)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				"kafkaClient.FetchMessage(ctx): %s", err.Error(),
			)
		}

		rpcMsg := kafkaMsgToRPC(msg)

		if err := stream.Send(rpcMsg); err != nil {
			return status.Errorf(
				codes.Internal,
				"stream.Send(rpcMsg): %s", err.Error(),
			)
		}

		kafkaClient.CommitMessages(ctx, msg)
	}
}

func kafkaMsgToRPC(msg kafka.Message) *pbV1.SubscribeTopicResponse {
	return &pbV1.SubscribeTopicResponse{
		Key:       string(msg.Key),
		Value:     string(msg.Value),
		Topic:     msg.Topic,
		Timestamp: uint64(msg.Time.Unix()),
	}

}
