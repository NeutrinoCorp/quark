package pkg

import (
	"strconv"

	"github.com/Shopify/sarama"
)

//	Implements sarama.ConsumerGroupHandler
type defaultKafkaConsumer struct {
	worker *kafkaWorker
}

func (k defaultKafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (k defaultKafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (k *defaultKafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Note: DO NOT SEND ANY ERRORS BACK IF YOU DONT WANT TO STOP THE CONSUMER GROUP'S SESSION (All workers)
	for msgConsumer := range claim.Messages() {
		eventCtx := session.Context()
		h := PopulateKafkaEventHeaders(msgConsumer)
		h.Set(HeaderKafkaMemberId, session.MemberID())
		h.Set(HeaderKafkaGenerationId, strconv.Itoa(int(session.GenerationID())))
		h.Set(HeaderConsumerGroup, k.worker.parent.Consumer.group)
		e := &Event{
			Context:    eventCtx,
			Topic:      msgConsumer.Topic,
			Header:     h,
			Body:       UnmarshalKafkaHeaders(msgConsumer),
			RawSession: session,
		}

		if k.worker.parent.Consumer.handler != nil {
			if commit := k.worker.parent.Consumer.handler.ServeEvent(k.worker.parent.setDefaultEventWriter(), e); commit {
				session.MarkMessage(msgConsumer, "")
			}
		}
		if k.worker.parent.Consumer.handlerFunc != nil {
			if commit := k.worker.parent.Consumer.handlerFunc(k.worker.parent.setDefaultEventWriter(), e); commit {
				session.MarkMessage(msgConsumer, "")
			}
		}

	}
	return nil
}
