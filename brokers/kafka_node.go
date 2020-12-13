package brokers

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
	"github.com/neutrinocorp/quark/internal"
)

type kafkaNode struct {
	Cluster  []string
	Config   *sarama.Config
	Consumer quark.Consumer
	Wg       *sync.WaitGroup
}

func (k kafkaNode) run(ctx context.Context, errChan chan error) {
	go func() {
		err := k.startConsumer(ctx, k.Consumer.Workers, k.Consumer.Topic, k.Consumer.Queue, k.Consumer.Handler,
			k.Consumer.ErrorHandler)
		if err != nil {
			if k.Consumer.ErrorHandler != nil {
				k.Consumer.ErrorHandler(ctx, k.Consumer.Topic, k.Consumer.Queue, err)
			}
			errChan <- err
			return
		}
	}()
	// retry consumers, ignore if not defined
	if k.Consumer.Retry != nil {
		go func() {
			err := k.startConsumer(ctx, k.Consumer.Retry.Workers, k.Consumer.Retry.Topic, k.Consumer.Retry.Queue,
				k.Consumer.Retry.Handler, k.Consumer.Retry.ErrorHandler)
			if err != nil {
				if k.Consumer.Retry.ErrorHandler != nil {
					k.Consumer.Retry.ErrorHandler(ctx, k.Consumer.Retry.Topic, k.Consumer.Retry.Queue, err)
				}
				errChan <- err
				return
			}
		}()
	}
	// DLQ consumers, ignore if not defined
	if k.Consumer.DeadLetter != nil {
		go func() {
			err := k.startConsumer(ctx, k.Consumer.DeadLetter.Workers, k.Consumer.DeadLetter.Topic,
				k.Consumer.DeadLetter.Queue, k.Consumer.DeadLetter.Handler, k.Consumer.DeadLetter.ErrorHandler)
			if err != nil {
				if k.Consumer.DeadLetter.ErrorHandler != nil {
					k.Consumer.DeadLetter.ErrorHandler(ctx, k.Consumer.DeadLetter.Topic, k.Consumer.DeadLetter.Queue,
						err)
				}
				errChan <- err
				return
			}
		}()
	}
}

func (k kafkaNode) startConsumer(ctx context.Context, workers uint, topic, queue string, handler interface{},
	errorHandler quark.ErrorHandler) error {
	defer k.Wg.Done()

	h, err := internal.ParseSaramaConsumer(handler)
	if err != nil {
		return err
	}

	// Kafka Consumer groups behaves like a queue, a queue is required by each topic
	consumer, err := sarama.NewConsumerGroup(k.Cluster, queue, k.Config)
	if err != nil {
		return err
	}
	defer func() {
		err = consumer.Close()
	}()

	if k.Config.Consumer.Return.Errors {
		go func() {
			if err := <-consumer.Errors(); err != nil {
				errorHandler(ctx, topic, queue, err)
			}
		}()
	}

	return k.forkWorkers(ctx, workers, consumer, topic, h)
}

func (k kafkaNode) forkWorkers(ctx context.Context, numWorkers uint, consumer sarama.ConsumerGroup, topic string,
	handler sarama.ConsumerGroupHandler) error {
	workerErrChan := make(chan error, numWorkers)
	for i := 0; i < int(numWorkers); i++ {
		go func() {
			workerErrChan <- k.consume(ctx, consumer, topic, handler)
		}()
	}

	return <-workerErrChan
}

func (k kafkaNode) consume(ctx context.Context, consumer sarama.ConsumerGroup, topic string,
	handler sarama.ConsumerGroupHandler) error {
	return consumer.Consume(ctx, []string{topic}, handler)
}
