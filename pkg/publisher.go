package pkg

import (
	"context"
)

// Publisher pushes the given Message into the Event-Driven ecosystem.
type Publisher interface {
	Publish(context.Context, *Message) error
}
