package pkg

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
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
	//	Remember Quark uses Message's "Kind" field as topic name, so the developer must specify it
	WriteMessage(context.Context, ...*Message) error
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

func (d *defaultEventWriter) Write(ctx context.Context, msg []byte, topics ...string) error {
	if d.publisher == nil {
		return ErrPublisherNotImplemented
	} else if len(topics) == 0 {
		return ErrNotEnoughTopics
	}
	errs := new(multierror.Error)
	for _, t := range topics {
		log.Printf(t)
		m := NewMessage(t, msg)
		d.parseHeader(m)
		if m.Metadata.RedeliveryCount >= d.node.setDefaultMaxRetries() {
			continue // avoid distributed loops (at macro scale)
		}
		time.Sleep(d.node.setDefaultRetryBackoff() * time.Duration(m.Metadata.RedeliveryCount))
		if err := d.publisher.Publish(ctx, m); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
	}
	return errs.ErrorOrNil()
}

// TODO: Add new return value (int) to tell user how many messages were published/written
func (d *defaultEventWriter) WriteMessage(ctx context.Context, msgs ...*Message) error {
	if d.publisher == nil {
		return ErrPublisherNotImplemented
	}
	errs := new(multierror.Error)
	for _, msg := range msgs {
		d.parseHeader(msg)
		if msg.Metadata.RedeliveryCount >= d.node.setDefaultMaxRetries() {
			continue // avoid distributed loops (at macro scale)
		}
		time.Sleep(d.node.setDefaultRetryBackoff() * time.Duration(msg.Metadata.RedeliveryCount))
		if err := d.publisher.Publish(ctx, msg); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
	}
	return errs.ErrorOrNil()
}

func (d *defaultEventWriter) parseHeader(msg *Message) {
	for k, v := range d.header {
		switch k {
		case HeaderMessageKind:
			msg.Kind = v
		case HeaderMessageCorrelationId:
			msg.Metadata.CorrelationId = v
		case HeaderMessageRedeliveryCount:
			if c, err := strconv.Atoi(v); err == nil {
				msg.Metadata.RedeliveryCount = c
			}
		case HeaderMessageHost:
			msg.Metadata.Host = v
		default:
			msg.Metadata.ExternalData[k] = v
		}
	}
}
