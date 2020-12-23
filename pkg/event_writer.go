package pkg

import (
	"context"
	"strconv"
	"time"
)

// EventWriter works as an Event response writer.
// Lets an Event to respond, fan-out a topic message or to fail properly by sending the failed Event into either a
// Retry or Dead Letter Queue (DLQ).
//
// Uses a Publisher to write the actual message
type EventWriter interface {
	// Publisher actual publisher used to push Event(s)
	Publisher() Publisher
	// Header Event metadata
	Header() Header
	// Write push the given encoded message into the Event-Driven ecosystem.
	//
	// Returns non-nil error if publisher failed to push Event or
	// returns ErrNotEnoughTopics if no topic was specified
	Write(ctx context.Context, msg []byte, topics ...string) error
	// WriteMessage push the given message into the Event-Driven ecosystem.
	//
	// Returns non-nil error if publisher failed to push Event
	//	This function is limited to push the given message to only one topic (specified as "Kind" in Message fields)
	WriteMessage(context.Context, *Message) error
}

type defaultEventWriter struct {
	node      *node
	publisher Publisher
	header    Header
}

// newEventWriter allocates and creates a default EventWriter
func newEventWriter(n *node, p Publisher) EventWriter {
	return &defaultEventWriter{
		node:      n,
		publisher: p,
		header:    Header{},
	}
}

func (d defaultEventWriter) Publisher() Publisher {
	return d.publisher
}

func (d *defaultEventWriter) Header() Header {
	return d.header
}

func (d defaultEventWriter) Write(ctx context.Context, msg []byte, topics ...string) error {
	if d.publisher == nil {
		return ErrPublisherNotImplemented
	} else if len(topics) == 0 {
		return ErrNotEnoughTopics
	}
	for _, t := range topics {
		m := NewMessage(t, msg)
		d.parseHeader(m)
		if m.Metadata.RedeliveryCount >= d.node.setDefaultMaxRetries() {
			return nil // avoid loops
		}
		time.Sleep(d.node.setDefaultRetryBackoff() * time.Duration(m.Metadata.RedeliveryCount))
		if err := d.publisher.Publish(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

func (d defaultEventWriter) WriteMessage(ctx context.Context, msg *Message) error {
	if d.publisher == nil {
		return ErrPublisherNotImplemented
	}
	d.parseHeader(msg)
	if msg.Metadata.RedeliveryCount >= d.node.setDefaultMaxRetries() {
		return nil // avoid loops
	}
	time.Sleep(d.node.setDefaultRetryBackoff() * time.Duration(msg.Metadata.RedeliveryCount))
	return d.publisher.Publish(ctx, msg)
}

func (d defaultEventWriter) parseHeader(msg *Message) {
	for k, v := range d.header {
		switch k {
		case "topic":
			msg.Kind = v
		case "correlation_id":
			msg.Metadata.CorrelationId = v
		case "redelivery_count":
			if c, err := strconv.Atoi(v); err == nil {
				msg.Metadata.RedeliveryCount = c
			}
		case "host":
			msg.Metadata.Host = v
		default:
			msg.Metadata.ExternalData[k] = v
		}
	}
}
