package quark

import (
	"context"

	"github.com/Shopify/sarama"
)

// KafkaConfiguration Apache Kafka specific Broker and Consumer configuration, overrides default values,
// contains from basic configuration to functions serving as Hooks when an action was dispatched
type KafkaConfiguration struct {
	Config   *sarama.Config
	Consumer KafkaConsumerConfig
	Producer KafkaProducerConfig
}

// KafkaConsumerConfig Apache Kafka consumer configuration
type KafkaConsumerConfig struct {
	GroupHandler     sarama.ConsumerGroupHandler
	PartitionHandler KafkaPartitionConsumer
	Topic            KafkaConsumerTopicConfig
	// Hooks
	OnReceived func(context.Context, *sarama.ConsumerMessage)
}

// KafkaConsumerTopicConfig Apache Kafka configuration used to override default consuming values
type KafkaConsumerTopicConfig struct {
	Partition int32
	Offset    int64
}

// KafkaProducerConfig Apache Kafka producer configuration
type KafkaProducerConfig struct {
	// Hooks
	OnSent func(ctx context.Context, message *sarama.ProducerMessage, partition int32, offset int64)
}
