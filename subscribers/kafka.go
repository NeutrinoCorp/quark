package subscribers

import (
	"log"

	"github.com/Shopify/sarama"
)

//	Implements sarama.ConsumerGroupHandler
type KafkaConsumerHandler struct{}

func (k KafkaConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (k KafkaConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (k KafkaConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msgConsumer := range claim.Messages() {
		log.Print(msgConsumer.Topic)
		log.Print(string(msgConsumer.Value))
		session.MarkMessage(msgConsumer, "")
	}
	return nil
}
