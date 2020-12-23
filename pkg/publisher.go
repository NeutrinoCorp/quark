package pkg

import (
	"context"
)

// Publisher pushes the given Message into the Event-Driven ecosystem.
type Publisher interface {
	Publish(context.Context, ...*Message) error
}

func getDefaultPublisher(provider string, providerCfg interface{}, cluster []string) Publisher {
	// broker and node already ensured providers, still this is a safe casting and if publisher is nil, event writer
	// will still return an error when trying to write messages
	switch provider {
	case KafkaProvider:
		if cfg, ok := providerCfg.(KafkaConfiguration); ok {
			return &defaultKafkaPublisher{
				cfg:     cfg,
				cluster: cluster,
			}
		}
	}
	return nil
}
