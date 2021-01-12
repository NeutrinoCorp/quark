package quark

const (
	// KafkaProvider Quark for Apache Kafka
	KafkaProvider = "kafka"
	// AWSProvider Quark for Amazon Web Services (SNS, SQS, SNS+SQS, Event Bridge)
	AWSProvider = "aws"
	// GCPProvider Quark for Google Cloud Platform (PubSub)
	GCPProvider = "gcp"
	// RabbitMqProvider Quark for RabbitMQ
	RabbitMqProvider = "rabbitmq"
	// ActiveMqProvider Quark for Apache ActiveMQ
	ActiveMqProvider = "activemq"
	// NatsProvider Quark for NATS
	NatsProvider = "nats"
)

func ensureValidProvider(provider string, providerCfg interface{}) error {
	switch provider {
	case KafkaProvider:
		if _, ok := providerCfg.(KafkaConfiguration); ok {
			return nil
		}
	case AWSProvider:
		if _, ok := providerCfg.(AWSConfiguration); ok {
			return nil
		}
	}

	return ErrProviderNotValid
}
