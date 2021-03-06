package quark

import (
	"time"
)

// Message an Event's body.
//
// It is intended to be passed into the actual message broker/queue system (e.g. Kafka, AWS SNS/SQS, RabbitMQ) as
// default message form.
//
// Thus, it keeps struct consistency between messages published in any programming language.
//
// It is based on the Async API specification.
type Message struct {
	// Id message unique identifier
	Id string `json:"message_id"`
	// Kind message topic, it is recommended to use the Async API topic naming convention
	Kind string `json:"kind"`
	// PublishTime time message was pushed into the ecosystem
	PublishTime time.Time `json:"publish_time"`
	// Attributes actual event encoded in binary (e.g. JSON or Apache Avro)
	Attributes []byte `json:"attributes"`
	// Metadata message volatile information
	Metadata metadata `json:"metadata"`
}

type metadata struct {
	// CorrelationId root message id
	CorrelationId string `json:"correlation_id"`
	// Host sender node ip address
	Host string `json:"host"`
	// RedeliveryCount attempts this specific message tried to get process
	RedeliveryCount int `json:"redelivery_count"`
	// ExternalData non-Quark data may be stored here (e.g. non-Quark headers)
	ExternalData map[string]string `json:"external_data"`
}

// NewMessage creates a new message without a parent message
func NewMessage(id, kind string, attributes []byte) *Message {
	return &Message{
		Id:          id,
		Kind:        kind,
		PublishTime: time.Now().UTC(),
		Attributes:  attributes,
		Metadata: metadata{
			CorrelationId:   id,
			Host:            "",
			RedeliveryCount: 0,
			ExternalData:    map[string]string{},
		},
	}
}

// NewMessageFromParent creates a new Message setting as parent (on the CorrelationId field) the given parentId
func NewMessageFromParent(parentId, id, kind string, attributes []byte) *Message {
	return &Message{
		Id:          id,
		Kind:        kind,
		PublishTime: time.Now().UTC(),
		Attributes:  attributes,
		Metadata: metadata{
			CorrelationId:   parentId,
			Host:            "",
			RedeliveryCount: 0,
			ExternalData:    map[string]string{},
		},
	}
}

// Encode returns the encoded attributes
func (m Message) Encode() ([]byte, error) {
	return m.Attributes, nil
}

// Length bytes size of the current message's attributes
func (m Message) Length() int {
	return len(m.Attributes)
}
