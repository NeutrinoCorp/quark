package quark

import (
	"strconv"

	"github.com/Shopify/sarama"
)

func PopulateKafkaEventHeaders(msg *sarama.ConsumerMessage) Header {
	h := Header{}
	h.Set(HeaderKafkaPartition, strconv.Itoa(int(msg.Partition)))
	h.Set(HeaderKafkaOffset, strconv.Itoa(int(msg.Offset)))
	h.Set(HeaderKafkaKey, string(msg.Key))
	h.Set(HeaderKafkaValue, string(msg.Value))
	h.Set(HeaderKafkaTimestamp, msg.Timestamp.String())
	h.Set(HeaderKafkaBlockTimestamp, msg.BlockTimestamp.String())
	for _, f := range msg.Headers {
		h.Set(string(f.Key), string(f.Value))
	}

	return h
}
