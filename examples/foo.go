package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/neutrinocorp/quark"
)

type AWSPublisher struct{}

func (a AWSPublisher) Publish(ctx context.Context, msgs ...*quark.Message) error {
	log.Print("publishing message to AWS")
	return nil
}

type NotificationHandler struct{}

func (h NotificationHandler) ServeEvent(_ quark.EventWriter, e *quark.Event) bool {
	log.Print("received notification from " + e.Topic)
	log.Print(e.Body)
	return true
}

func LogErrors() func(context.Context, error) {
	return func(ctx context.Context, err error) {
		log.Print(err)
	}
}

func main() {
	// BDD clause
	// Create broker
	b, err := quark.NewBroker(quark.KafkaProvider, quark.KafkaConfiguration{
		Config: nil,
	})
	if err != nil {
		panic(err)
	}
	b.Cluster = []string{"localhost:9092"}

	// Example: Chat, communication between multiple topics
	b.Topic("chat.0").PoolSize(4).HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
		_, _ = w.Write(e.Context, []byte("hello"), "chat.1")
		return true
	})
	b.Topic("chat.1").HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
		_, _ = w.Write(e.Context, []byte("goodbye"), "chat.0")
		return true
	})

	// Example: Listen to multiple notifications using specific resiliency configurations
	b.Topics("bob.notifications", "alice.notifications").MaxRetries(5).RetryBackoff(time.Second * 3).
		Handle(NotificationHandler{})

	// Example: Listen to some user trading using custom publisher provider and sending a response to multiple topics
	b.Topic("alex.trades").Publisher(AWSPublisher{}).HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
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
	b.Topic("retry.truck.0.gps").Provider(quark.KafkaProvider).Address("localhost:9092", "localhost:9093").
		HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
			if e.Body.Metadata.RedeliveryCount >= 3 {
				return true // avoid loops
			}
			e.Body.Metadata.RedeliveryCount++
			w.Header().Set(quark.HeaderMessageRedeliveryCount, strconv.Itoa(e.Body.Metadata.RedeliveryCount))
			msg := quark.NewMessageFromParent(e.Body.Metadata.CorrelationId, e.Body.Kind, e.Body.Attributes)
			_, _ = w.WriteMessage(e.Context, msg)
			// _ = w.Publisher().Publish(e.Context, e.Body) is also valid but will not write given headers
			return true
		})

	b.ErrorHandler = LogErrors()

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
