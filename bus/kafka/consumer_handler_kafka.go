package kafka

import (
	"context"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
)

// KafkaPartitionConsumer This consumer is the default way to consume messages from Apache Kafka.
//
// It pulls messages from an specific partition inside an Apache Kafka cluster or Broker. This way of consuming messages
// is useful when actual parallelization of the process itself is required. When in a pool, it will pull messages for
// each consumer running in a Worker, running the process at the same time in a worker pool.
type KafkaPartitionConsumer interface {
	Consume(context.Context, sarama.PartitionConsumer, *quark.Consumer, quark.EventWriter)
}

type defaultKafkaPartitionConsumer struct {
	worker *kafkaWorker
}

// Consume starts consuming from a single Apache Kafka partition
func (k *defaultKafkaPartitionConsumer) Consume(ctx context.Context, p sarama.PartitionConsumer, c *quark.Consumer,
	e quark.EventWriter) {
	for msgConsumer := range p.Messages() {
		eventCtx := ctx
		if k.worker.cfg.Consumer.OnReceived != nil {
			k.worker.cfg.Consumer.OnReceived(eventCtx, msgConsumer)
		}
		h := NewKafkaHeader(msgConsumer)
		h.Set(HeaderHighWaterMarkOffset, strconv.Itoa(int(p.HighWaterMarkOffset())))
		// set up required parent data (tracing, redelivery and correlation)
		e.ReplaceHeader(newQuarkHeaders(h))
		body := new(quark.Message)
		UnmarshalKafkaMessage(msgConsumer, body)
		ev := &quark.Event{
			Context:    eventCtx,
			Topic:      msgConsumer.Topic,
			Header:     h,
			Body:       body,
			RawValue:   msgConsumer.Value,
			RawSession: p,
		}

		if handler := c.GetHandle(); handler != nil {
			handler.ServeEvent(e, ev)
		}
		if handleFunc := c.GetHandleFunc(); handleFunc != nil {
			handleFunc(e, ev)
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
		h := NewKafkaHeader(msgConsumer)
		h.Set(HeaderMemberId, session.MemberID())
		h.Set(HeaderGenerationId, strconv.Itoa(int(session.GenerationID())))
		h.Set(quark.HeaderConsumerGroup, k.worker.parent.Consumer.GetGroup())
		body := new(quark.Message)
		UnmarshalKafkaMessage(msgConsumer, body)
		e := &quark.Event{
			Context:    eventCtx,
			Topic:      msgConsumer.Topic,
			Header:     h,
			Body:       body,
			RawValue:   msgConsumer.Value,
			RawSession: session,
		}

		if handler := k.worker.parent.Consumer.GetHandle(); handler != nil {
			// set up required parent data (tracing, redelivery and correlation)
			evWriter := k.worker.parent.GetEventWriter()
			evWriter.ReplaceHeader(newQuarkHeaders(h))
			if commit := handler.ServeEvent(evWriter, e); commit {
				session.MarkMessage(msgConsumer, "")
			}
		}
		if handlerFunc := k.worker.parent.Consumer.GetHandleFunc(); handlerFunc != nil {
			// set up required parent data (tracing, redelivery and correlation)
			evWriter := k.worker.parent.GetEventWriter()
			evWriter.ReplaceHeader(newQuarkHeaders(h))
			if commit := handlerFunc(evWriter, e); commit {
				session.MarkMessage(msgConsumer, "")
			}
		}

	}
	return nil
}

func newQuarkHeaders(h quark.Header) quark.Header {
	hEv := quark.Header{}
	hEv.Set(quark.HeaderSpanContext, h.Get(quark.HeaderSpanContext))
	hEv.Set(quark.HeaderMessageCorrelationId, h.Get(quark.HeaderMessageCorrelationId))
	hEv.Set(quark.HeaderMessageRedeliveryCount, h.Get(quark.HeaderMessageRedeliveryCount))
	return hEv
}
