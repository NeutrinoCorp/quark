package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
)

// NewKafkaBroker allocates and returns a Kafka Broker
func NewKafkaBroker(cfg *sarama.Config, opts ...quark.Option) *quark.Broker {
	broker := quark.NewBroker(opts...)
	setKafkaBrokerDefaults(cfg, broker)
	return broker
}

func setKafkaBrokerDefaults(cfg *sarama.Config, b *quark.Broker) {
	setDefaultKafkaConfig(cfg, b)
	setDefaultKafkaPublisher(b)
	setDefaultKafkaWorkerFactory(b)
}

func setDefaultKafkaConfig(cfg *sarama.Config, b *quark.Broker) {
	if b.ProviderConfig == nil {
		b.ProviderConfig = KafkaConfiguration{
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
	}
}

func setDefaultKafkaPublisher(b *quark.Broker) {
	if kafkaCfg, ok := b.ProviderConfig.(KafkaConfiguration); ok && b.Publisher == nil {
		b.Publisher = &KafkaPublisher{
			cfg:     kafkaCfg,
			cluster: b.Cluster,
		}
	}
}

func setDefaultKafkaWorkerFactory(b *quark.Broker) {
	if kafkaCfg, ok := b.ProviderConfig.(KafkaConfiguration); ok && b.WorkerFactory == nil {
		b.WorkerFactory = newKafkaWorkerFactory(kafkaCfg)
	}
}

func newKafkaWorkerFactory(cfg KafkaConfiguration) quark.WorkerFactory {
	return func(parent *quark.Supervisor) quark.Worker {
		return &kafkaWorker{
			id:     0,
			parent: parent,
			cfg:    cfg,
		}
	}
}
