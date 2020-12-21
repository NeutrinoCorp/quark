package pkg

import (
	"context"
	"strconv"
)

// EventWriter works as an Event response writer.
// Lets an Event to respond or fail properly by sending the failed Event into either a Retry or Dead Letter Queue (DLQ).
//
// Uses a Publisher to write the actual message
type EventWriter interface {
	// Publisher actual publisher used to push Event(s)
	Publisher() Publisher
	// Header Event metadata
	Header() Header
	// Write push the given encoded message into the Event-Driven ecosystem
	Write(ctx context.Context, topic string, msg []byte) error
	// WriteMessage push the given message into the Event-Driven ecosystem
	WriteMessage(context.Context, *Message) error
}

type defaultEventWriter struct {
	publisher Publisher
	header    Header
}

// NewEventWriter allocates and creates a default EventWriter
func NewEventWriter(p Publisher) EventWriter {
	return &defaultEventWriter{
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

func (d defaultEventWriter) Write(ctx context.Context, topic string, msg []byte) error {
	if d.publisher == nil {
		return ErrPublisherNotFound
	}
	m := NewMessage(topic, msg)
	d.parseHeader(m)
	return d.publisher.Publish(ctx, m)
}

func (d defaultEventWriter) WriteMessage(ctx context.Context, msg *Message) error {
	if d.publisher == nil {
		return ErrPublisherNotFound
	}
	d.parseHeader(msg)
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
