package kafka

const (
	// HeaderPartition Topic partition where Message was stored in Apache Kafka commit log
	HeaderPartition = "quark-kafka-partition"
	// HeaderOffset Topic partition offset, item number inside an specific Topic partition
	HeaderOffset = "quark-kafka-offset"
	// HeaderKey Message key on a Apache Kafka Topic, used to group multiple messages by its ID, it can be also used
	// to compact commit logs by updating and then deleting previous item versions (removes duplicates)
	HeaderKey = "quark-kafka-key"
	// HeaderValue Apache Kafka message body in binary format
	HeaderValue = "quark-kafka-value"
	// HeaderTimestamp Apache Kafka insertion time
	HeaderTimestamp = "quark-kafka-timestamp"
	// HeaderBlockTimestamp Apache Kafka producer insertion time
	HeaderBlockTimestamp = "quark-kafka-block-timestamp"
	// HeaderMemberId Member unique identifier from an Apache Kafka Consumer Group
	HeaderMemberId = "quark-kafka-member-id"
	// HeaderGenerationId Unique identifier when Topic partition list is requested when joining an Apache Kafka
	// Consumer Group
	HeaderGenerationId = "quark-kafka-generation-id"
	// HeaderHighWaterMarkOffset Last message that was successfully copied to all of the log’s replicas in an Apache
	// Kafka cluster
	HeaderHighWaterMarkOffset = "quark-kafka-high-water-mark-offset"
)
