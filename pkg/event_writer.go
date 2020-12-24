package pkg

import (
	"context"
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
	injectHeader(Header)
	// Publisher actual publisher used to push Event(s)
	Publisher() Publisher
	// Header Event metadata
	Header() Header
	// Write push the given encoded message into the Event-Driven ecosystem.
	//
	// Returns number of messages published and non-nil error if publisher failed to push Event or
	// returns ErrNotEnoughTopics if no topic was specified
	//	Sometimes, the writer might not publish messages to broker since they have passed the maximum redelivery cap
	Write(ctx context.Context, msg []byte, topics ...string) (int, error)
	// WriteMessage push the given message into the Event-Driven ecosystem.
	//
	// Returns number of messages published
	// and non-nil error if publisher failed to push Event
	//	Remember Quark uses Message's "Kind" field as topic name, so the developer must specify it either in mentioned field or in response headers
	//	Sometimes, the writer might not publish messages to broker since they have passed the maximum redelivery cap
	WriteMessage(context.Context, ...*Message) (int, error)
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

func (d *defaultEventWriter) injectHeader(h Header) {
	d.header = h
}

func (d defaultEventWriter) Publisher() Publisher {
	return d.publisher
}

func (d *defaultEventWriter) Header() Header {
	return d.header
}

func (d *defaultEventWriter) Write(ctx context.Context, msg []byte, topics ...string) (int, error) {
	if d.publisher == nil {
		return 0, ErrPublisherNotImplemented
	} else if len(topics) == 0 {
		return 0, ErrNotEnoughTopics
	}
	errs := new(multierror.Error)
	msgPublished := 0
	for _, t := range topics {
		m := NewMessage(t, msg)
		d.parseHeader(m)
		if m.Metadata.RedeliveryCount > d.node.setDefaultMaxRetries() {
			continue // avoid distributed loops (at macro scale)
		}
		time.Sleep(d.node.setDefaultRetryBackoff() * time.Duration(m.Metadata.RedeliveryCount))
		if err := d.publisher.Publish(ctx, m); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		msgPublished++
	}
	return msgPublished, errs.ErrorOrNil()
}

func (d *defaultEventWriter) WriteMessage(ctx context.Context, msgs ...*Message) (int, error) {
	if d.publisher == nil {
		return 0, ErrPublisherNotImplemented
	} else if len(msgs) == 0 {
		return 0, ErrNotEnoughTopics
	}
	errs := new(multierror.Error)
	msgPublished := 0
	for _, msg := range msgs {
		d.parseHeader(msg)
		if msg.Metadata.RedeliveryCount > d.node.setDefaultMaxRetries() {
			continue // avoid distributed loops (at macro scale)
		}
		time.Sleep(d.node.setDefaultRetryBackoff() * time.Duration(msg.Metadata.RedeliveryCount))
		if err := d.publisher.Publish(ctx, msg); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		msgPublished++
	}
	return msgPublished, errs.ErrorOrNil()
}

func (d *defaultEventWriter) parseHeader(msg *Message) {
	for k, v := range d.header {
		switch k {
		case HeaderMessageKind:
			msg.Kind = v
		case HeaderMessageCorrelationId:
			if v != "" {
				msg.Metadata.CorrelationId = v
			}
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
