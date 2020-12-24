package quark

const (
	HeaderMessageId              = "quark-message-id"
	HeaderMessageKind            = "quark-kind"
	HeaderMessagePublishTime     = "quark-publish-time"
	HeaderMessageAttributes      = "quark-attributes"
	HeaderMessageCorrelationId   = "quark-metadata-correlation-id"
	HeaderMessageHost            = "quark-metadata-host"
	HeaderMessageRedeliveryCount = "quark-metadata-redelivery-count"

	HeaderConsumerGroup = "quark-consumer-group"
	HeaderSpanContext   = "quark-span-context"

	HeaderKafkaPartition           = "quark-kafka-partition"
	HeaderKafkaOffset              = "quark-kafka-offset"
	HeaderKafkaKey                 = "quark-kafka-key"
	HeaderKafkaValue               = "quark-kafka-value"
	HeaderKafkaTimestamp           = "quark-kafka-timestamp"
	HeaderKafkaBlockTimestamp      = "quark-kafka-block-timestamp"
	HeaderKafkaMemberId            = "quark-kafka-member-id"
	HeaderKafkaGenerationId        = "quark-kafka-generation-id"
	HeaderKafkaHighWaterMarkOffset = "quark-kafka-high-water-mark-offset"
)
