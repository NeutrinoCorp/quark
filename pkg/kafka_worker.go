package pkg

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"

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

// might delete this since we already have default groups using topics names
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

	k.setDefaultConsumerPartitionHandler().Consume(ctx, k.partitioner, k.parent.Consumer,
		k.parent.setDefaultEventWriter())

	return nil
}

func (k *kafkaWorker) Close() error {
	log.Print(k.parent.Consumer.topics, "closing worker")
	errs := new(multierror.Error)
	if k.group != nil {
		if err := k.group.Close(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	if k.consumer != nil {
		if err := k.consumer.Close(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	if k.partitioner != nil {
		if err := k.partitioner.Close(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs.ErrorOrNil()
}

func (k *kafkaWorker) setDefaultConsumerGroupHandler() sarama.ConsumerGroupHandler {
	if k.cfg.ConsumerGroupHandler != nil {
		return k.cfg.ConsumerGroupHandler
	}

	return &defaultKafkaConsumer{
		worker: k,
	}
}

func (k *kafkaWorker) setDefaultConsumerPartitionHandler() KafkaPartitionConsumer {
	if k.cfg.ConsumerPartitionHandler != nil {
		return k.cfg.ConsumerPartitionHandler
	}
	return &defaultKafkaPartitionConsumer{}
}
