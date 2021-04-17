package quark

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/jpillora/backoff"
)

// EventWriter works as an Event response writer.
// Lets an Event to respond, fan-out a topic message or to fail properly by sending the failed Event into either a
// Retry or Dead Letter Queue (DLQ).
//
// Uses a Publisher to write the actual message
type EventWriter interface {
	ReplaceHeader(Header)
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
	//	Remember Quark uses Message's "Type" field as topic name, so the developer must specify it either in mentioned field or in response headers
	//	Sometimes, the writer might not publish messages to broker since they have passed the maximum redelivery cap
	WriteMessage(context.Context, ...*Message) (int, error)
	// WriteRetry push the given Event into a retry topic.
	// Publishing aside, WriteRetry method adds a new ID for every message passed using the Broker's ID Generator/Factory.
	//
	// It is recommended to point the Message's Topic/Type to an specific retry topic/queue (e.g. foo.executed -> foo.executed.retry/foo.executed.retry.N).
	//
	// This implementation differs from others because it increments the given Message "redelivery_count" delta field by one
	WriteRetry(ctx context.Context, msg *Message) error
}

// ErrMessageRedeliveredTooMuch the message has been published the number of times of the configuration limit
var ErrMessageRedeliveredTooMuch = errors.New("message has been redelivered too much")

type defaultEventWriter struct {
	Node      *Node
	publisher Publisher
	header    Header
	backoff   *backoff.Backoff
}

// newEventWriter allocates and creates a default EventWriter
func newEventWriter(n *Node, p Publisher) EventWriter {
	return &defaultEventWriter{
		Node:      n,
		publisher: p,
		header:    Header{},
		backoff: &backoff.Backoff{
			Factor: 2,
			Jitter: true,
			Min:    n.setDefaultRetryBackoff(),
			Max:    n.setDefaultRetryBackoff() * 3,
		},
	}
}

func (d *defaultEventWriter) ReplaceHeader(h Header) {
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
		m := NewMessage(d.Node.Broker.setDefaultMessageIdGenerator()(), t, msg)
		if err := d.publish(ctx, m); err != nil {
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
		if err := d.publish(ctx, msg); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		msgPublished++
	}
	return msgPublished, errs.ErrorOrNil()
}

func (d *defaultEventWriter) WriteRetry(ctx context.Context, msg *Message) error {
	if d.publisher == nil {
		return ErrPublisherNotImplemented
	} else if msg == nil {
		return ErrEmptyMessage
	}

	msg.Id = d.Node.Broker.setDefaultMessageIdGenerator()()
	msg.Metadata.RedeliveryCount++
	d.Header().Set(HeaderMessageRedeliveryCount, strconv.Itoa(msg.Metadata.RedeliveryCount))
	return d.publish(ctx, msg)
}

func (d *defaultEventWriter) publish(ctx context.Context, msg *Message) error {
	d.marshalMessage(msg)

	backoffFactor := msg.Metadata.RedeliveryCount
	if backoffFactor > d.Node.setDefaultMaxRetries() {
		return ErrMessageRedeliveredTooMuch
	}

	time.Sleep(d.backoff.ForAttempt(float64(backoffFactor)))
	return d.publisher.Publish(ctx, msg)
}

func (d *defaultEventWriter) marshalMessage(msg *Message) {
	for k, v := range d.header {
		switch k {
		case HeaderMessageType:
			msg.Type = v
		case HeaderMessageCorrelationId:
			if v != "" {
				msg.Metadata.CorrelationId = v
			}
		case HeaderMessageRedeliveryCount:
			if c, err := strconv.Atoi(v); msg.Type == d.header.Get(HeaderMessageType) && err == nil {
				msg.Metadata.RedeliveryCount = c
			}
		case HeaderMessageHost:
			msg.Metadata.Host = v
		default:
			msg.Metadata.ExternalData[k] = v
		}
	}

	if d.Node != nil {
		msg.Source = d.Node.setDefaultSource()
		msg.ContentType = d.Node.setDefaultContentType()
	}
}
