package kafka

import (
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
)

// NewKafkaHeader creates a Message Header from an Apache Kafka message
func NewKafkaHeader(msg *sarama.ConsumerMessage) quark.Header {
	h := quark.Header{}
	h.Set(HeaderPartition, strconv.Itoa(int(msg.Partition)))
	h.Set(HeaderOffset, strconv.Itoa(int(msg.Offset)))
	h.Set(HeaderKey, string(msg.Key))
	h.Set(HeaderValue, string(msg.Value))
	h.Set(HeaderTimestamp, msg.Timestamp.String())
	h.Set(HeaderBlockTimestamp, msg.BlockTimestamp.String())
	for _, f := range msg.Headers {
		h.Set(string(f.Key), string(f.Value))
	}

	return h
}
