package quark

// using CNCF CloudEvents specification v1.0
// ref. https://cloudevents.io/
// ref. https://github.com/cloudevents/spec/blob/v1.0.1/spec.md
const (
	// HeaderMessageId Message ID
	HeaderMessageId = "quark-id"
	// HeaderMessageType This attribute contains a value describing the type of event related to the originating occurrence.
	HeaderMessageType = "quark-type"
	// HeaderMessageSpecVersion The version of the CloudEvents specification which the event uses.
	HeaderMessageSpecVersion = "quark-spec-version"
	// HeaderMessageSource Identifies the context in which an event happened.
	HeaderMessageSource = "quark-source"
	// HeaderMessageTime Time when Message was published
	HeaderMessageTime = "quark-time"
	// HeaderMessageDataContentType Content type of data value. This attribute enables data to carry any type of content,
	// whereby format and encoding might differ from that of the chosen event format.
	HeaderMessageDataContentType = "quark-content-type"
	// HeaderMessageDataSchema Identifies the schema that data adheres to.
	HeaderMessageDataSchema = "quark-data-schema"
	// HeaderMessageData Message body encoded bytes
	HeaderMessageData = "quark-data"
	// HeaderMessageSubject This describes the subject of the event in the context of the event producer (identified by source).
	HeaderMessageSubject = "quark-subject"
	// HeaderMessageCorrelationId Message parent (origin)
	HeaderMessageCorrelationId = "quark-metadata-correlation-id"
	// HeaderMessageHost Node IP from deployed cluster in infrastructure
	HeaderMessageHost = "quark-metadata-host"
	// HeaderMessageRedeliveryCount Message total redeliveries
	HeaderMessageRedeliveryCount = "quark-metadata-redelivery-count"
	// HeaderMessageError Message error message from processing pipeline
	HeaderMessageError = "quark-metadata-error"

	// HeaderConsumerGroup Consumer group this message was received by
	HeaderConsumerGroup = "quark-consumer-group"
	// HeaderSpanContext Message span parent, used for distributed tracing mechanisms such as OpenCensus, OpenTracing
	// and/or OpenTelemetry
	HeaderSpanContext = "quark-span-context"
)
