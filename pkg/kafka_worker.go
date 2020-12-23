package pkg

import (
	"context"
	"log"
	"strconv"

	"github.com/Shopify/sarama"
)

type kafkaWorker struct {
	parent *node
	cfg    KafkaConfiguration

	group       sarama.ConsumerGroup
	consumer    sarama.Consumer
	partitioner sarama.PartitionConsumer
}

func (k *kafkaWorker) Parent() *node {
	return k.parent
}

func (k *kafkaWorker) StartJob(ctx context.Context) error {
	log.Print(k.parent.Consumer.topics, k.parent.setDefaultProvider())
	if err := k.ensureGroup(); err != nil {
		return err
	} else if len(k.parent.Consumer.topics) > 1 || k.parent.Consumer.group != "" {
		return k.startConsumerGroup(ctx)
	}
	return k.startConsumer(ctx)
}

func (k *kafkaWorker) ensureGroup() error {
	if len(k.parent.Consumer.topics) > 1 && k.parent.Consumer.group == "" {
		return ErrRequiredGroup
	}

	return nil
}

func (k *kafkaWorker) startConsumerGroup(ctx context.Context) error {
	group, err := sarama.NewConsumerGroup(k.parent.setDefaultCluster(), k.parent.setDefaultGroup(), k.cfg.Config)
	if err != nil {
		return err
	}
	k.group = group

	if k.cfg.Config.Consumer.Return.Errors && k.parent.Broker.ErrorHandler != nil {
		go k.parent.Broker.ErrorHandler(<-k.group.Errors())
	}

	for {
		if err = k.group.Consume(ctx, k.parent.Consumer.topics, k.setDefaultConsumerGroupHandler()); err != nil {
			return err
		}
	}
}

func (k *kafkaWorker) startConsumer(ctx context.Context) error {
	consumer, err := sarama.NewConsumer(k.parent.setDefaultCluster(), k.cfg.Config)
	if err != nil {
		return err
	}
	k.consumer = consumer

	cPartition, err := k.consumer.ConsumePartition(k.parent.Consumer.topics[0], k.cfg.ConsumerTopic.Partition,
		k.cfg.ConsumerTopic.Offset)
	if err != nil {
		return err
	}
	k.partitioner = cPartition

	if k.cfg.Config.Consumer.Return.Errors && k.parent.Broker.ErrorHandler != nil {
		go k.parent.Broker.ErrorHandler(<-k.partitioner.Errors())
	}

	for msgConsumer := range k.partitioner.Messages() {
		eventCtx := ctx
		h := PopulateKafkaEventHeaders(msgConsumer)
		h.Set(HeaderKafkaHighWaterMarkOffset, strconv.Itoa(int(k.partitioner.HighWaterMarkOffset())))
		e := &Event{
			Context:    eventCtx,
			Topic:      msgConsumer.Topic,
			Header:     h,
			Body:       UnmarshalKafkaHeaders(msgConsumer),
			RawSession: k.partitioner,
		}

		if k.parent.Consumer.handler != nil {
			k.parent.Consumer.handler.ServeEvent(k.parent.setDefaultEventWriter(), e)
		}
		if k.parent.Consumer.handlerFunc != nil {
			k.parent.Consumer.handlerFunc(k.parent.setDefaultEventWriter(), e)
		}
	}

	return nil
}

func (k *kafkaWorker) Close() error {
	log.Print(k.parent.Consumer.topics, "closing worker")
	if k.group != nil {
		if err := k.group.Close(); err != nil {
			return err
		}
	}
	if k.consumer != nil {
		if err := k.consumer.Close(); err != nil {
			return err
		}
	}
	if k.partitioner != nil {
		if err := k.partitioner.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (k *kafkaWorker) setDefaultConsumerGroupHandler() sarama.ConsumerGroupHandler {
	if k.cfg.ConsumerGroupHandler != nil {
		return k.cfg.ConsumerGroupHandler
	}

	return &defaultKafkaConsumer{
		worker: k,
	}
}
