package kafka

import (
	"context"
	"errors"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/go-multierror"
	"github.com/neutrinocorp/quark"
)

type kafkaWorker struct {
	id     int
	parent *quark.Supervisor
	cfg    KafkaConfiguration

	group       sarama.ConsumerGroup
	consumer    sarama.Consumer
	partitioner sarama.PartitionConsumer
}

func (k *kafkaWorker) SetID(i int) {
	k.id = i
}

func (k *kafkaWorker) Parent() *quark.Supervisor {
	return k.parent
}

func (k *kafkaWorker) StartJob(ctx context.Context) error {
	if err := k.ensureGroup(); err != nil {
		return err
	} else if len(k.parent.Consumer.GetTopics()) > 1 || k.parent.Consumer.GetGroup() != "" {
		return k.startConsumerGroup(ctx)
	}
	return k.startConsumer(ctx)
}

// might delete this since we already have default groups using topics names
func (k *kafkaWorker) ensureGroup() error {
	if len(k.parent.Consumer.GetTopics()) > 1 && k.parent.Consumer.GetGroup() == "" {
		return quark.ErrRequiredGroup
	}

	return nil
}

func (k *kafkaWorker) startConsumerGroup(ctx context.Context) error {
	group, err := sarama.NewConsumerGroup(k.parent.GetCluster(), k.parent.GetGroup(), k.cfg.Config)
	if err != nil {
		return err
	}
	k.group = group

	if k.cfg.Config.Consumer.Return.Errors && k.parent.Broker.ErrorHandler != nil {
		go func() {
			for e := range k.group.Errors() {
				if k.parent.Broker.ErrorHandler != nil {
					k.parent.Broker.ErrorHandler(ctx, e)
				}
			}
		}()
	}

	// blocking I/O
	go func() {
		retries := 0
		for {
			err = k.group.Consume(ctx, k.parent.Consumer.GetTopics(), k.setDefaultConsumerGroupHandler())
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return
			} else if errors.Is(err, sarama.ErrOutOfBrokers) {
				if retries <= k.parent.Broker.GetConnRetries() {
					retries++
					time.Sleep(k.parent.Broker.GetConnRetryBackoff() * time.Duration(retries))
					continue
				}
			}
			if err != nil {
				if k.parent.Broker.ErrorHandler != nil {
					go k.parent.Broker.ErrorHandler(ctx, err)
				}
				return
			}
			retries = 0 // restart count since problem was fixed
		}
	}()

	return nil
}

func (k *kafkaWorker) startConsumer(ctx context.Context) error {
	consumer, err := sarama.NewConsumer(k.parent.GetCluster(), k.cfg.Config)
	if err != nil {
		return err
	}
	k.consumer = consumer

	cPartition, err := k.consumer.ConsumePartition(k.parent.Consumer.GetTopics()[0], k.cfg.Consumer.Topic.Partition,
		k.cfg.Consumer.Topic.Offset)
	if err != nil {
		return err
	}
	k.partitioner = cPartition

	if k.cfg.Config.Consumer.Return.Errors && k.parent.Broker.ErrorHandler != nil {
		go func() {
			for e := range k.partitioner.Errors() {
				if k.parent.Broker.ErrorHandler != nil {
					k.parent.Broker.ErrorHandler(ctx, e)
				}
			}
		}()
	}

	// Blocking I/O
	go func() {
		k.setDefaultConsumerPartitionHandler().Consume(ctx, k.partitioner, k.parent.Consumer,
			k.parent.GetEventWriter())
	}()

	return nil
}

func (k *kafkaWorker) Close() error {
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
	if k.cfg.Consumer.GroupHandler != nil {
		return k.cfg.Consumer.GroupHandler
	}

	return &defaultKafkaConsumer{
		worker: k,
	}
}

func (k *kafkaWorker) setDefaultConsumerPartitionHandler() KafkaPartitionConsumer {
	if k.cfg.Consumer.PartitionHandler != nil {
		return k.cfg.Consumer.PartitionHandler
	}
	return &defaultKafkaPartitionConsumer{
		worker: k,
	}
}
