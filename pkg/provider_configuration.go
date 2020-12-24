package pkg

import (
	"context"

	"github.com/Shopify/sarama"
)

type KafkaConfiguration struct {
	Config   *sarama.Config
	Consumer kafkaConsumerConfig
	Producer kafkaProducerConfig
}

type kafkaConsumerConfig struct {
	GroupHandler     sarama.ConsumerGroupHandler
	PartitionHandler KafkaPartitionConsumer
	Topic            kafkaConsumerTopicConfig
	// Hooks
	OnReceived func(context.Context, *sarama.ConsumerMessage)
}

type kafkaConsumerTopicConfig struct {
	Partition int32
	Offset    int64
}

type kafkaProducerConfig struct {
	// Hooks
	OnSent func(ctx context.Context, message *sarama.ProducerMessage, partition int32, offset int64)
}
