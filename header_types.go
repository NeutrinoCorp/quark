package quark

const (
	// HeaderMessageId Message ID
	HeaderMessageId = "quark-message-id"
	// HeaderMessageKind Message type (command or domain event)
	HeaderMessageKind = "quark-kind"
	// HeaderMessagePublishTime Time when Message was published
	HeaderMessagePublishTime = "quark-publish-time"
	// HeaderMessageAttributes Message body encoded bytes
	HeaderMessageAttributes = "quark-attributes"
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

	// HeaderKafkaPartition Topic partition where Message was stored in Apache Kafka commit log
	HeaderKafkaPartition = "quark-kafka-partition"
	// HeaderKafkaOffset Topic partition offset, item number inside an specific Topic partition
	HeaderKafkaOffset = "quark-kafka-offset"
	// HeaderKafkaKey Message key on a Apache Kafka Topic, used to group multiple messages by its ID, it can be also used
	// to compact commit logs by updating and then deleting previous item versions (removes duplicates)
	HeaderKafkaKey = "quark-kafka-key"
	// HeaderKafkaValue Apache Kafka message body in binary format
	HeaderKafkaValue = "quark-kafka-value"
	// HeaderKafkaTimestamp Apache Kafka insertion time
	HeaderKafkaTimestamp = "quark-kafka-timestamp"
	// HeaderKafkaBlockTimestamp Apache Kafka producer insertion time
	HeaderKafkaBlockTimestamp = "quark-kafka-block-timestamp"
	// HeaderKafkaMemberId Member unique identifier from an Apache Kafka Consumer Group
	HeaderKafkaMemberId = "quark-kafka-member-id"
	// HeaderKafkaGenerationId Unique identifier when Topic partition list is requested when joining an Apache Kafka
	// Consumer Group
	HeaderKafkaGenerationId = "quark-kafka-generation-id"
	// HeaderKafkaHighWaterMarkOffset Last message that was successfully copied to all of the log’s replicas in an Apache
	// Kafka cluster
	HeaderKafkaHighWaterMarkOffset = "quark-kafka-high-water-mark-offset"
)
