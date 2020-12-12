package brokers

import (
	"context"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/go-multierror"
	"github.com/neutrinocorp/quark"
	"github.com/neutrinocorp/quark/internal"
)

type Kafka struct {
	Cluster      []string
	Config       *sarama.Config
	ErrorHandler quark.ErrorHandler

	consumers map[string]quark.Consumer
	mu        *sync.Mutex
}

func NewKafka(cluster []string, config *sarama.Config, errHandler quark.ErrorHandler) *Kafka {
	return &Kafka{
		Cluster:      cluster,
		Config:       config,
		ErrorHandler: errHandler,
		consumers:    map[string]quark.Consumer{},
		mu:           new(sync.Mutex),
	}
}

func (k *Kafka) Register(consumers ...quark.Consumer) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if err := k.ensureConsumerUniqueness(consumers...); err != nil {
		return err
	}

	for _, consumer := range consumers {
		k.consumers[consumer.Topic] = consumer
	}
	return nil
}

func (k Kafka) ensureConsumerUniqueness(consumers ...quark.Consumer) error {
	errs := new(multierror.Error)
	for _, consumer := range consumers {
		if k.IsConsumerRegistered(consumer.Topic) {
			errs = multierror.Append(errs, quark.ErrConsumerAlreadyRegistered)
		}
	}
	return errs.ErrorOrNil()
}

func (k Kafka) IsConsumerRegistered(name string) bool {
	_, ok := k.consumers[name]
	return ok
}

func (k Kafka) ListenAndServe(ctx context.Context) error {
	var err error
	for {
		if err = k.startWorkers(ctx); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return err
		}
	}
}

func (k Kafka) startWorkers(ctx context.Context) error {
	errs := make(chan error)
	for _, consumer := range k.consumers {
		consumer := consumer
		// main subscribers
		go func() {
			err := k.startConsumer(ctx, consumer.Workers, consumer.Topic, consumer.Queue, consumer.Handler)
			if err != nil {
				k.ErrorHandler(ctx, consumer.Topic, consumer.Queue, err)
				errs <- err
				return
			}
		}()
		// retry subscribers
		go func() {
			err := k.startConsumer(ctx, consumer.Retry.Workers, consumer.Retry.Topic, consumer.Retry.Queue,
				consumer.Retry.Handler)
			if err != nil {
				k.ErrorHandler(ctx, consumer.Retry.Topic, consumer.Retry.Queue, err)
				errs <- err
				return
			}
		}()
		// DLQ subscribers
		go func() {
			err := k.startConsumer(ctx, consumer.DeadLetter.Workers, consumer.DeadLetter.Topic,
				consumer.DeadLetter.Queue, consumer.DeadLetter.Handler)
			if err != nil {
				k.ErrorHandler(ctx, consumer.DeadLetter.Topic, consumer.DeadLetter.Queue, err)
				errs <- err
				return
			}
		}()
	}
	select {
	case err := <-errs:
		return err
	}
}

func (k Kafka) startConsumer(ctx context.Context, workers int, topic, queue string, handler interface{}) error {
	h, err := internal.ParseSaramaConsumer(handler)
	if err != nil {
		return err
	}

	return k.forkWorkers(ctx, workers, topic, queue, h)
}

func (k Kafka) forkWorkers(ctx context.Context, numWorkers int, name, queue string,
	handler sarama.ConsumerGroupHandler) error {
	errChan := make(chan error)
	for i := 0; i < numWorkers; i++ {
		go func() {
			if err := k.consume(ctx, name, queue, handler); err != nil {
				errChan <- err
				return
			}
			errChan <- nil
		}()
	}

	select {
	case err := <-errChan:
		return err
	}
}

func (k Kafka) consume(ctx context.Context, name, queue string, handler sarama.ConsumerGroupHandler) error {
	consumerG, err := sarama.NewConsumerGroup(k.Cluster, queue, k.Config)
	if err != nil {
		return err
	}
	defer func() {
		err = consumerG.Close()
		log.Print("closing ", name)
	}()

	log.Print("consuming ", name, " with queue ", queue)

	ctxConsumer, cancel := context.WithCancel(ctx)
	defer cancel()
	return consumerG.Consume(ctxConsumer, []string{name}, handler)
}
