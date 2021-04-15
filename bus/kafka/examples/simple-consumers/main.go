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

func main() {
	// Create broker
	b := kafka.NewKafkaBroker(newSaramaCfg(), "localhost:9092")

	b.ErrorHandler = func(ctx context.Context, err error) {
		log.Print(err)
	}

	b.Topics("chat.0").PoolSize(1).MaxRetries(3).
		HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
			log.Printf("topic: %s | message: %s", e.Topic, e.RawValue)
			log.Printf("topic: %s | redelivery: %s", e.Topic, e.Header.Get(quark.HeaderMessageRedeliveryCount))
			_, _ = w.Write(e.Context, []byte(e.Header.Get(kafka.HeaderKafkaValue)), "chat.1")
			return true
		})

	b.Topics("chat.1").PoolSize(1).MaxRetries(5).
		HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
			redelivery := e.Header.Get(quark.HeaderMessageRedeliveryCount)
			log.Printf("topic: %s | message: %s", e.Topic, e.RawValue)
			log.Printf("topic: %s | redelivery: %s", e.Topic, redelivery)
			_, _ = w.Write(e.Context, []byte(e.Header.Get(kafka.HeaderKafkaValue)), "chat.1")
			return true
		})

	// graceful shutdown
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	go func() {
		if err := b.ListenAndServe(); err != nil && err != quark.ErrBrokerClosed {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Printf("stopping %d nodes and %d workers", b.RunningNodes(), b.RunningWorkers())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := b.Shutdown(ctx); err != nil {
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
