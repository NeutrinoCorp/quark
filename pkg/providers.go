package pkg

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
	}

	return ErrProviderNotValid
}
