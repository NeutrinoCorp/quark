package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
	"github.com/neutrinocorp/quark/bus/kafka"
)

var customErrHandler = func(ctx context.Context, err error) {
	log.Print(err)
}

func main() {
	// Create broker
	b := kafka.NewKafkaBroker(
		newSaramaCfg(),
		quark.WithCluster("localhost:19092", "localhost:29092", "localhost:39092"),
		quark.WithBaseMessageSource("https://neutrinocorp.org/cloudevents"),
		quark.WithBaseMessageContentType("message/partial"),
		quark.WithErrorHandler(customErrHandler))

	b.Topics("org.neutrinocorp.command.user.create").Group("user-service").PoolSize(5).MaxRetries(3).
		HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
			log.Printf("topic: %s | message: %s", e.Topic, e.RawValue)
			log.Printf("topic: %s | redelivery: %s", e.Topic, e.Header.Get(quark.HeaderMessageRedeliveryCount))
			_, _ = w.Write(e.Context, []byte(e.Body.Data), "org.neutrinocorp.event.user.created")
			return true
		})

	b.Topics("org.neutrinocorp.event.user.created").Group("notification-service").PoolSize(5).MaxRetries(3).
		HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
			redelivery := e.Header.Get(quark.HeaderMessageRedeliveryCount)
			log.Printf("topic: %s | message: %s", e.Topic, e.RawValue)
			log.Printf("topic: %s | redelivery: %s", e.Topic, redelivery)

			/*
				// Retry strategy
				err := w.WriteRetry(e.Context, e.Body)
				if errors.Is(err, quark.ErrMessageRedeliveredTooMuch) {
					_, _ = w.Write(e.Context, e.Body.Data, "dlq.chat.1")
				}*/
			return true
		})

	b.Topics("org.neutrinocorp.event.user.created").Group("analytics-service").PoolSize(5).MaxRetries(3).
		HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
			redelivery := e.Header.Get(quark.HeaderMessageRedeliveryCount)
			log.Printf("topic: %s | message: %s", e.Topic, e.RawValue)
			log.Printf("topic: %s | redelivery: %s", e.Topic, redelivery)
			_, _ = w.Write(e.Context, e.Body.Data, "org.neutrinocorp.event.user.created.dlq")
			return true
		})

	b.Topic("org.neutrinocorp.event.user.created.dlq").Group("foo-group").
		HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
			log.Printf("topic: %s | message: %s", e.Topic, string(e.Body.Data))
			return true
		})

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		if err := b.ListenAndServe(); err != nil && err != quark.ErrBrokerClosed {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Printf("stopping %d supervisor(s) and %d worker(s)", b.ActiveSupervisors(), b.ActiveWorkers())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := b.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Print(b.ActiveSupervisors(), b.ActiveWorkers()) // should be 0,0
}

func newSaramaCfg() *sarama.Config {
	config := sarama.NewConfig()
	config.ClientID = "neutrino-sample"
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Consumer.Retry.Backoff = time.Second * 10
	return config
}
