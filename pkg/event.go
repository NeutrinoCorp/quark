package pkg

import (
	"context"
)

// Event represents something that has happened or requested to happen within the Event-Driven ecosystem.
//
// For example, this could be either a domain event or an asynchronous command.
type Event struct {
	Context context.Context
	Topic   string
	Header  Header
	Body    *Message
}
