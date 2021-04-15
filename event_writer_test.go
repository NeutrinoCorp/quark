package quark

import (
	"context"
	"errors"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var eventWriterHeaderTestingSuite = []struct {
	in  Header
	exp *Message
}{
	{
		in: Header{
			HeaderMessageKind:            "foo",
			HeaderMessageCorrelationId:   "123",
			HeaderMessageHost:            "192.168.1.1",
			HeaderMessageRedeliveryCount: "1",
			HeaderMessageError:           "cassandra: foo bar error",
		},
		exp: &Message{
			Id:          "",
			Kind:        "foo",
			PublishTime: time.Time{},
			Attributes:  nil,
			Metadata: metadata{
				CorrelationId:   "123",
				Host:            "192.168.1.1",
				RedeliveryCount: 1,
				ExternalData: map[string]string{
					HeaderMessageError: "cassandra: foo bar error",
				},
			},
		},
	},
}

func TestEventWriterHeader(t *testing.T) {
	for _, tt := range eventWriterHeaderTestingSuite {
		t.Run("Event Writer header manipulation", func(t *testing.T) {
			w := &defaultEventWriter{
				node:      nil,
				publisher: nil,
				header:    tt.in,
			}
			msg := &Message{
				Id:          "",
				Kind:        "",
				PublishTime: time.Time{},
				Attributes:  nil,
				Metadata: metadata{
					CorrelationId:   "",
					Host:            "",
					RedeliveryCount: 0,
					ExternalData:    map[string]string{},
				},
			}
			w.marshalMessage(msg)
			assert.Exactly(t, tt.exp, msg)
			w.injectHeader(nil) // overrides
			assert.Nil(t, w.Header())
		})
	}
}

var eventWriterPublishTestingSuite = []struct {
	publisher  Publisher
	topics     []string
	redelivery int
	exp        error
	expWritten int
}{
	{&stubPublisher{fail: false}, []string{}, 0, ErrNotEnoughTopics, 0},
	{nil, []string{"foo"}, 0, ErrPublisherNotImplemented, 0},
	{&stubPublisher{fail: true}, []string{"foo"}, 0, errStubPublisher, 0},
	{&stubPublisher{fail: false}, []string{"foo"}, 5, nil, 0},
	{&stubPublisher{fail: false}, []string{"foo"}, 0, nil, 1},
	{&stubPublisher{fail: false}, []string{"foo"}, 4, nil, 1},
	{&stubPublisher{fail: false}, []string{"foo", "bar", "baz"}, 0, nil, 3},
}

func TestEventWriterWrite(t *testing.T) {
	for _, tt := range eventWriterPublishTestingSuite {
		ctx := context.Background()
		t.Run("Event Writer write", func(t *testing.T) {
			w := newEventWriter(&node{Consumer: &Consumer{}, Broker: &Broker{
				MaxRetries:   5,
				RetryBackoff: time.Millisecond * 150,
			}}, tt.publisher)
			w.Header().Set(HeaderMessageRedeliveryCount, strconv.Itoa(tt.redelivery))
			if tt.publisher == nil {
				assert.Nil(t, w.Publisher())
			} else {
				assert.NotNil(t, w.Publisher())
			}
			m, err := w.Write(ctx, []byte("You're a rockstar"), tt.topics...)
			log.Print(err)
			assert.True(t, errors.Is(err, tt.exp))
			assert.Equal(t, tt.expWritten, m)
		})
	}
}

func TestDefaultEventWriter_WriteMessage(t *testing.T) {
	for _, tt := range eventWriterPublishTestingSuite {
		ctx := context.Background()
		t.Run("Event Writer write message", func(t *testing.T) {
			w := newEventWriter(&node{Consumer: &Consumer{}, Broker: &Broker{
				MaxRetries:   5,
				RetryBackoff: time.Millisecond * 150,
			}}, tt.publisher)
			w.Header().Set(HeaderMessageRedeliveryCount, strconv.Itoa(tt.redelivery))
			if tt.publisher == nil {
				assert.Nil(t, w.Publisher())
			} else {
				assert.NotNil(t, w.Publisher())
			}
			msgs := make([]*Message, 0)
			for _, topic := range tt.topics {
				msgs = append(msgs, &Message{
					Id:          "123",
					Kind:        topic,
					PublishTime: time.Time{},
					Attributes:  nil,
					Metadata: metadata{
						CorrelationId:   "123",
						Host:            "192.168.1.1",
						RedeliveryCount: 0,
						ExternalData: map[string]string{
							HeaderMessageError: "cassandra: foo bar error",
						},
					},
				})
			}
			m, err := w.WriteMessage(ctx, msgs...)
			assert.True(t, errors.Is(err, tt.exp))
			assert.Equal(t, tt.expWritten, m)
		})
	}
}
