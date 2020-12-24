package quark

import (
	"context"

	"github.com/Shopify/sarama"
)

type KafkaConfiguration struct {
	Config   *sarama.Config
	Consumer KafkaConsumerConfig
	Producer KafkaProducerConfig
}

type KafkaConsumerConfig struct {
	GroupHandler     sarama.ConsumerGroupHandler
	PartitionHandler KafkaPartitionConsumer
	Topic            KafkaConsumerTopicConfig
	// Hooks
	OnReceived func(context.Context, *sarama.ConsumerMessage)
}

type KafkaConsumerTopicConfig struct {
	Partition int32
	Offset    int64
}

type KafkaProducerConfig struct {
	// Hooks
	OnSent func(ctx context.Context, message *sarama.ProducerMessage, partition int32, offset int64)
}
