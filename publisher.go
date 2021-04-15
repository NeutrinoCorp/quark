package quark

import (
	"context"
)

// Publisher pushes the given Message into the Event-Driven ecosystem.
type Publisher interface {
	Publish(context.Context, ...*Message) error
}

// PublisherFactory is a crucial Broker and/or Consumer component which generates the concrete publishers
// Quark will use to produce data
type PublisherFactory func(providerCfg interface{}, cluster []string) Publisher
