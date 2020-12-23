package pkg

import (
	"context"
)

// Publisher pushes the given Message into the Event-Driven ecosystem.
type Publisher interface {
	Publish(context.Context, ...*Message) error
}

func getDefaultPublisher(provider string) Publisher {
	switch provider {
	case KafkaProvider:
		return &defaultKafkaPublisher{}
	default:
		return nil
	}
}
