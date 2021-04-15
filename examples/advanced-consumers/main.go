package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
)

type awsPublisherStub struct{}

func (a awsPublisherStub) Publish(ctx context.Context, msgs ...*quark.Message) error {
	for _, msg := range msgs {
		log.Printf("publishing - message: %s", msg.Kind)
	}
	return nil
}

type notificationHandlerStub struct{}

func (h notificationHandlerStub) ServeEvent(_ quark.EventWriter, e *quark.Event) bool {
	log.Printf("topic: %s | message: %s", e.Topic, e.RawValue)
	log.Printf("topic: %s | correlation: %s", e.Topic, e.Header.Get(quark.HeaderMessageCorrelationId))
	return true
}

func logErrors() func(context.Context, error) {
	return func(ctx context.Context, err error) {
		log.Print(err)
	}
}

func main() {
	// BDD clause
	// Create broker
	b, err := quark.NewBroker(quark.KafkaProvider, quark.KafkaConfiguration{
		Config: newSaramaCfg(),
		Consumer: quark.KafkaConsumerConfig{
			GroupHandler:     nil,
			PartitionHandler: nil,
			Topic: quark.KafkaConsumerTopicConfig{
				Partition: 0,
				Offset:    sarama.OffsetNewest,
			},
			OnReceived: nil,
		},
	})
	if err != nil {
		panic(err)
	}
	b.Cluster = []string{"localhost:19092", "localhost:29092", "localhost:39092"}

	// Example: Listen to multiple notifications using specific resiliency configurations
	b.Topics("bob.notifications", "alice.notifications").Group("notifications").MaxRetries(5).RetryBackoff(time.Second * 3).
		Handle(notificationHandlerStub{})

	// Example: Listen to some user trading using custom publisher provider and sending a response to multiple topics
	b.Topic("alex.trades").Group("alex.trades").Publisher(awsPublisherStub{}).
		HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
			_, _ = w.Write(e.Context, []byte("alex has traded in a new index fund"),
				"aws.alex.trades", "aws.analytics.trades")
			return true
		})

	// Example: Listen to a feed failing completely (send message to DLQ)
	b.Topic("alice.feed").PoolSize(10).HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
		_, _ = w.Write(e.Context, []byte("failed to process message"), "dlq.feed")
		return true
	})

	// Example: Truck GPS tracker using a custom provider and address, fail temporarily (sending message to retry queue)
	b.Topic("retry.truck.0.gps").Group("retry.truck.0.gps").Provider(quark.KafkaProvider).MaxRetries(3).
		Address("localhost:9092", "localhost:9093").RetryBackoff(time.Second * 3).
		HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
			log.Printf("topic: %s | message: %s", e.Topic, e.RawValue)
			log.Printf("topic: %s | redelivery: %d", e.Topic, e.Body.Metadata.RedeliveryCount)
			_, _ = w.Write(e.Context, e.RawValue, e.Topic)
			// _ = w.Publisher().Publish(e.Context, e.Body) is also valid but will not write given headers
			return true
		})

	b.ErrorHandler = logErrors()

	// graceful shutdown
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	go func() {
		if err = b.ListenAndServe(); err != nil && err != quark.ErrBrokerClosed {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Printf("stopping %d nodes and %d workers", b.RunningNodes(), b.RunningWorkers())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err = b.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Print(b.RunningNodes(), b.RunningWorkers()) // should be 0,0
}

func newSaramaCfg() *sarama.Config {
	config := sarama.NewConfig()
	config.ClientID = "neutrino-sample"
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	return config
}
