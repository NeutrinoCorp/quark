package pkg

import "github.com/Shopify/sarama"

type KafkaConfiguration struct {
	Config               *sarama.Config
	ConsumerGroupHandler sarama.ConsumerGroupHandler
	ConsumerTopic        kafkaConsumerTopicConfig
}

type kafkaConsumerTopicConfig struct {
	Partition int32
	Offset    int64
}

type AWSConfiguration struct{}
