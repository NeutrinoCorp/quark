package quark

const (
	KafkaProvider    = "kafka"
	AWSProvider      = "aws"
	ActiveMqProvider = "activemq"
	RabbitMqProvider = "rabbitmq"
	NatsProvider     = "nats"
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
