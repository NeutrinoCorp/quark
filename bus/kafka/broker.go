package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
)

// NewKafkaBroker allocates and returns a Kafka Broker
func NewKafkaBroker(cfg *sarama.Config, addrs ...string) *quark.Broker {
	kafkaCfg := KafkaConfiguration{
		Config: cfg,
		Consumer: KafkaConsumerConfig{
			GroupHandler:     nil,
			PartitionHandler: nil,
			Topic: KafkaConsumerTopicConfig{
				Partition: 0,
				Offset:    sarama.OffsetNewest,
			},
			OnReceived: nil,
		},
		Producer: KafkaProducerConfig{},
	}
	broker := quark.NewBroker(quark.KafkaProvider, kafkaCfg)
	broker.Cluster = addrs
	broker.Publisher = &KafkaPublisher{
		cfg:     kafkaCfg,
		cluster: addrs,
	}
	broker.WorkerFactory = newKafkaWorkerFactory(kafkaCfg)

	return broker
}

func newKafkaWorkerFactory(cfg KafkaConfiguration) quark.WorkerFactory {
	return func(parent *quark.Node) quark.Worker {
		return &kafkaWorker{
			id:     0,
			parent: parent,
			cfg:    cfg,
		}
	}
}
