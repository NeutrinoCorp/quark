package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark/pkg"
)

func main() {
	// Create broker
	b := pkg.NewKafkaBroker(newSaramaCfg(), "localhost:9092")

	b.ErrorHandler = func(ctx context.Context, err error) {
		log.Print(err)
	}

	// Example: Chat, communication between multiple topics
	b.Topics("chat.0", "chat.2").Group("neutrino-group-0").PoolSize(3).MaxRetries(3).
		HandleFunc(func(w pkg.EventWriter, e *pkg.Event) bool {
			log.Printf("topic: %s | message: %s", e.Topic, e.RawValue)
			log.Printf("topic: %s | redelivery: %s", e.Topic, e.Header.Get(pkg.HeaderMessageRedeliveryCount))
			w.Header().Set(pkg.HeaderMessageRedeliveryCount, strconv.Itoa(1))
			_, _ = w.Write(e.Context, []byte(e.Header.Get(pkg.HeaderKafkaValue)), "chat.1")
			return true
		})

	b.Topics("chat.1").Group("neutrino-group-1").PoolSize(3).MaxRetries(3).
		HandleFunc(func(w pkg.EventWriter, e *pkg.Event) bool {
			log.Printf("topic: %s | message: %s", e.Topic, e.RawValue)
			log.Printf("topic: %s | redelivery: %s", e.Topic, e.Header.Get(pkg.HeaderMessageRedeliveryCount))

			e.Body.Metadata.RedeliveryCount += 2
			w.Header().Set(pkg.HeaderMessageRedeliveryCount, strconv.Itoa(e.Body.Metadata.RedeliveryCount))

			_, _ = w.Write(e.Context, []byte("hello"), "chat.3")
			return true
		})

	b.Topic("chat.3").HandleFunc(func(w pkg.EventWriter, e *pkg.Event) bool {
		log.Printf("topic: %s | message: %s", e.Topic, e.RawValue)
		log.Printf("topic: %s | redelivery: %s", e.Topic, e.Header.Get(pkg.HeaderMessageRedeliveryCount))
		// _, _ = w.Write(e.Context, []byte("hello"), "chat.4", "chat.5")
		return true
	})

	// graceful shutdown
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	go func() {
		if err := b.ListenAndServe(); err != nil && err != pkg.ErrBrokerClosed {
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
