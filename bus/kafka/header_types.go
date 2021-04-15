package kafka

const (
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
