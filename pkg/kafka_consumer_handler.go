package pkg

import (
	"context"
	"strconv"

	"github.com/Shopify/sarama"
)

type KafkaPartitionConsumer interface {
	Consume(context.Context, sarama.PartitionConsumer, *Consumer, EventWriter)
}

type defaultKafkaPartitionConsumer struct {
	worker *kafkaWorker
}

func (k *defaultKafkaPartitionConsumer) Consume(ctx context.Context, p sarama.PartitionConsumer, c *Consumer,
	e EventWriter) {
	for msgConsumer := range p.Messages() {
		eventCtx := ctx
		if k.worker.cfg.Consumer.OnReceived != nil {
			k.worker.cfg.Consumer.OnReceived(eventCtx, msgConsumer)
		}
		h := PopulateKafkaEventHeaders(msgConsumer)
		h.Set(HeaderKafkaHighWaterMarkOffset, strconv.Itoa(int(p.HighWaterMarkOffset())))
		// set up required parent data (tracing, redelivery and correlation)
		hEv := Header{}
		hEv.Set(HeaderSpanContext, h.Get(HeaderSpanContext))
		hEv.Set(HeaderMessageCorrelationId, h.Get(HeaderMessageCorrelationId))
		hEv.Set(HeaderMessageRedeliveryCount, h.Get(HeaderMessageRedeliveryCount))
		e.injectHeader(hEv)
		ev := &Event{
			Context:    eventCtx,
			Topic:      msgConsumer.Topic,
			Header:     h,
			Body:       UnmarshalKafkaHeaders(msgConsumer),
			RawValue:   msgConsumer.Value,
			RawSession: p,
		}

		if c.handler != nil {
			c.handler.ServeEvent(e, ev)
		}
		if c.handlerFunc != nil {
			c.handlerFunc(e, ev)
		}
	}
}

//	Implements sarama.ConsumerGroupHandler
type defaultKafkaConsumer struct {
	worker *kafkaWorker
}

func (k *defaultKafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (k *defaultKafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (k *defaultKafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Note: DO NOT SEND ANY ERRORS BACK IF YOU DONT WANT TO STOP THE CONSUMER GROUP'S SESSION (All workers)
	for msgConsumer := range claim.Messages() {
		if k.worker.cfg.Consumer.OnReceived != nil {
			k.worker.cfg.Consumer.OnReceived(session.Context(), msgConsumer)
		}

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
			RawValue:   msgConsumer.Value,
			RawSession: session,
		}

		if k.worker.parent.Consumer.handler != nil {
			// set up required parent data (tracing, redelivery and correlation)
			evWriter := k.worker.parent.setDefaultEventWriter()
			hEv := Header{}
			hEv.Set(HeaderSpanContext, h.Get(HeaderSpanContext))
			hEv.Set(HeaderMessageCorrelationId, h.Get(HeaderMessageCorrelationId))
			hEv.Set(HeaderMessageRedeliveryCount, h.Get(HeaderMessageRedeliveryCount))
			evWriter.injectHeader(hEv)
			if commit := k.worker.parent.Consumer.handler.ServeEvent(evWriter, e); commit {
				session.MarkMessage(msgConsumer, "")
			}
		}
		if k.worker.parent.Consumer.handlerFunc != nil {
			// set up required parent data (tracing, redelivery and correlation)
			evWriter := k.worker.parent.setDefaultEventWriter()
			hEv := Header{}
			hEv.Set(HeaderSpanContext, h.Get(HeaderSpanContext))
			hEv.Set(HeaderMessageCorrelationId, h.Get(HeaderMessageCorrelationId))
			hEv.Set(HeaderMessageRedeliveryCount, h.Get(HeaderMessageRedeliveryCount))
			evWriter.injectHeader(hEv)
			if commit := k.worker.parent.Consumer.handlerFunc(evWriter, e); commit {
				session.MarkMessage(msgConsumer, "")
			}
		}

	}
	return nil
}
