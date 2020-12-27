package quark

import (
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalKafkaHeaders(t *testing.T) {
	t.Run("Unmarshal Kafka headers", func(t *testing.T) {
		msg := NewMessage("1", "chat.0", []byte("hello there"))

		publishTime, err := msg.PublishTime.MarshalBinary()
		if err != nil {
			publishTime = []byte(msg.PublishTime.String())
		}

		h := UnmarshalKafkaHeaders(&sarama.ConsumerMessage{
			Headers: []*sarama.RecordHeader{
				{[]byte(HeaderMessagePublishTime), publishTime},
			},
			Timestamp:      time.Time{},
			BlockTimestamp: time.Time{},
			Key:            nil,
			Value:          []byte("hello there"),
			Topic:          "",
			Partition:      0,
			Offset:         0,
		})

		assert.Equal(t, msg.PublishTime.String(), h.PublishTime.String())
	})
}
