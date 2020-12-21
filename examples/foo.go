package main

import (
	"log"
	"strconv"
	"time"

	"github.com/neutrinocorp/quark/pkg"
)

type NotificationHandler struct{}

func (h NotificationHandler) ServeEvent(_ pkg.EventWriter, e *pkg.Event) {
	log.Print("received notification from " + e.Topic)
	log.Print(e.Body)
}

func LogErrors() func(error) {
	return func(err error) {
		log.Print(err)
	}
}

func main() {
	// BDD clause
	// Create broker
	b := pkg.NewBroker()

	// Example: Chat
	b.Topic("chat.0").PoolSize(4).HandleFunc(func(w pkg.EventWriter, e *pkg.Event) {
		_ = w.Write(e.Context, "chat.1", []byte("hello"))
	})
	b.Topic("chat.1").HandleFunc(func(w pkg.EventWriter, e *pkg.Event) {
		_ = w.Write(e.Context, "chat.0", []byte("goodbye"))
	})

	// Example: Listen to multiple notifications
	b.Topics("bob.notifications", "alice.notifications").MaxRetries(5).RetryBackoff(time.Second * 3).
		Handle(NotificationHandler{})

	// Example: Listen to a feed, then fail
	b.Topic("alice.feed").PoolSize(10).HandleFunc(func(w pkg.EventWriter, e *pkg.Event) {
		_ = w.Write(e.Context, "dlq.feed", []byte("failed to process message"))
	})

	// Example: Truck GPS tracker, use custom provider and address, then fail temporarily (retry queue)
	b.Topic("retry.truck.0.gps").Provider(pkg.KafkaProvider).Address("localhost:9092", "localhost:9093").
		HandleFunc(func(w pkg.EventWriter, e *pkg.Event) {
			if e.Body.Metadata.RedeliveryCount >= 3 {
				return // avoid loops
			}

			msg := pkg.NewMessageFromParent(e.Body.Metadata.CorrelationId, e.Body.Kind, e.Body.Attributes)
			e.Body.Metadata.RedeliveryCount++
			w.Header().Set("redelivery_count", strconv.Itoa(e.Body.Metadata.RedeliveryCount))
			_ = w.WriteMessage(e.Context, msg)
			// _ = w.Publisher().Publish(e.Context, e.Body) is also valid but will not write given headers
		})

	b.ErrorHandler = LogErrors()

	_ = b.ListenAndServe()
}
