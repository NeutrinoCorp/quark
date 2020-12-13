package quark

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Message is the default unit of an event-driven architecture
//
// A Message will be passed to a Broker or external component to propagate or consume either a Domain Event or an
// Asynchronous Command
//	Implements sarama.Encoder
type Message struct {
	Id           string      `json:"message_id"`
	Kind         string      `json:"kind"`
	OccurredTime time.Time   `json:"occurred_time"`
	Attributes   interface{} `json:"attributes"`
	Metadata     metadata    `json:"metadata"`
}

type metadata struct {
	CorrelationId string      `json:"correlation_id"`
	Host          string      `json:"host"`
	SpanContext   interface{} `json:"span_context"`
}

func NewMessage(kind string, occurredTime time.Time, attributes interface{}) *Message {
	return &Message{
		Id:           uuid.New().String(),
		Kind:         kind,
		OccurredTime: occurredTime.UTC(),
		Attributes:   attributes,
		Metadata: metadata{
			CorrelationId: uuid.New().String(),
			Host:          getLocalIP(),
			SpanContext:   nil,
		},
	}
}

func NewMessageFromParent(parentId string, kind string, occurredTime time.Time, attributes interface{}) *Message {
	return &Message{
		Id:           uuid.New().String(),
		Kind:         kind,
		OccurredTime: occurredTime.UTC(),
		Attributes:   attributes,
		Metadata: metadata{
			CorrelationId: parentId,
			Host:          getLocalIP(),
			SpanContext:   nil,
		},
	}
}

func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnmarshalBinary(msg []byte) error {
	return json.Unmarshal(msg, m)
}

func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(*m)
}

func (m Message) Length() int {
	msg, err := m.Encode()
	if err != nil {
		return 0
	}
	return len(msg)
}