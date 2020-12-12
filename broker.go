package quark

import "context"

// Broker starts and controls Consumer(s)
type Broker interface {
	Register(consumers ...Consumer) error
	IsConsumerRegistered(topic string) bool
	ListenAndServe(context.Context) error
}
