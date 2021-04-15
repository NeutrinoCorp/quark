package quark

import (
	"time"
)

// CloudEventsVersion current Quark's CNCF CloudEvents specification version
const CloudEventsVersion = "1.0"

// Message an Event's body.
//
// It is intended to be passed into the actual message broker/queue system (e.g. Kafka, AWS SNS/SQS, RabbitMQ) as
// default message form.
//
// Thus, it keeps struct consistency between messages published in any programming language.
//
// Based on the CNCF CloudEvents specification v1.0
//
// ref. https://cloudevents.io/
//
// ref. https://github.com/cloudevents/spec/blob/v1.0.1/spec.md
type Message struct {
	// Id message unique identifier
	Id string `json:"id"`
	// Type message topic. This attribute contains a value describing the type of event related to the originating occurrence.
	//
	// Often this attribute is used for routing, observability, policy enforcement, etc.
	// The format of this is producer defined and might include information such as the version of the type -
	// see Versioning of Attributes in the Primer for more information.
	//
	//	e.g. com.example.object.deleted.v2
	Type string `json:"type"`
	// SpecVersion The version of the CloudEvents specification which the event uses.
	//
	// This enables the interpretation of the context.
	// Compliant event producers MUST use a value of 1.0 when referring to this version of the specification.
	SpecVersion string `json:"specversion"`
	// Source Identifies the context in which an event happened.
	// Often this will include information such as the type of the event source, the organization publishing the event or the process that produced the event.
	// The exact syntax and semantics behind the data encoded in the URI is defined by the event producer.
	//
	//	DNS e.g. /cloudevents/spec/pull/123
	//	Application-specific e.g. https://github.com/cloudevents
	Source string `json:"source"`

	// -- OPTIONAL --

	// Data The event payload. This specification does not place any restriction on the type of this information.
	//
	// It is encoded into a media format which is specified by the datacontenttype attribute (e.g. application/json),
	// and adheres to the dataschema format when those respective attributes are present.
	Data []byte `json:"data"`
	// ContentType Content type of data value. This attribute enables data to carry any type of content,
	// whereby format and encoding might differ from that of the chosen event format.
	ContentType string `json:"datacontenttype"`
	// DataSchema Identifies the schema that data adheres to.
	//
	// Incompatible changes to the schema SHOULD be reflected by a different URI. See Versioning of Attributes
	// in the Primer for more information
	DataSchema string `json:"dataschema"`
	// Subject This describes the subject of the event in the context of the event producer (identified by source).
	//
	// In publish-subscribe scenarios, a subscriber will typically subscribe to events emitted by a source, but the source
	// identifier alone might not be sufficient as a qualifier for any specific event if the source context has internal sub-structure.
	Subject string `json:"subject"`
	// Time Timestamp of when the occurrence happened.
	// If the time of the occurrence cannot be determined then this attribute MAY be set to some other time (such as the current time) by the
	// Quark producer, however all producers for the same source MUST be consistent in this respect.
	// In other words, either they all use the actual time of the occurrence or they all use the same algorithm to determine the value used.
	Time time.Time `json:"time"`

	// Metadata message volatile information
	Metadata MessageMetadata `json:"metadata"`
}

// MessageMetadata a message volatile fields, used to store useful Quark resiliency mechanisms and custom
// developer-defined fields
type MessageMetadata struct {
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
func NewMessage(id, msgType string, data []byte) *Message {
	return &Message{
		Id:          id,
		Type:        msgType,
		SpecVersion: CloudEventsVersion,
		Time:        time.Now().UTC(),
		Data:        data,
		Metadata: MessageMetadata{
			CorrelationId:   id,
			Host:            "",
			RedeliveryCount: 0,
			ExternalData:    map[string]string{},
		},
	}
}

// NewMessageFromParent creates a new Message setting as parent (on the CorrelationId field) the given parentId
func NewMessageFromParent(parentId, id, msgType string, data []byte) *Message {
	return &Message{
		Id:          id,
		Type:        msgType,
		SpecVersion: CloudEventsVersion,
		Time:        time.Now().UTC(),
		Data:        data,
		Metadata: MessageMetadata{
			CorrelationId:   parentId,
			Host:            "",
			RedeliveryCount: 0,
			ExternalData:    map[string]string{},
		},
	}
}

// Encode returns the encoded data
func (m Message) Encode() ([]byte, error) {
	return m.Data, nil
}

// Length bytes size of the current message's data
func (m Message) Length() int {
	return len(m.Data)
}
