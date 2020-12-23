package pkg

import (
	"context"

	"github.com/Shopify/sarama"
)

type KafkaConfiguration struct {
	Config                   *sarama.Config
	ConsumerGroupHandler     sarama.ConsumerGroupHandler
	ConsumerPartitionHandler KafkaPartitionConsumer
	ConsumerTopic            kafkaConsumerTopicConfig
	// Hook
	ProducerOnSent func(ctx context.Context, message *sarama.ProducerMessage, partition int32, offset int64)
}

type kafkaConsumerTopicConfig struct {
	Partition int32
	Offset    int64
}

type AWSConfiguration struct{}
