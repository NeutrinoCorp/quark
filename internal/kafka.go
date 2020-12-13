package internal

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
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
		msg := new(quark.Message)
		if err := msg.UnmarshalBinary(msgConsumer.Value); err != nil {
			// send to DLQ
			log.Printf("sending %s to DLQ", msgConsumer.Topic)
		}

		log.Printf("%+v", msg)
		session.MarkMessage(msgConsumer, "")
	}
	return nil
}
