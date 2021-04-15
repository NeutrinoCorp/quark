package kafka

import (
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalKafkaHeaders(t *testing.T) {
	t.Run("Unmarshal Kafka headers", func(t *testing.T) {
		msg := quark.NewMessage("1", "chat.0", []byte("hello there"))

		publishTime, err := msg.Time.MarshalBinary()
		if err != nil {
			publishTime = []byte(msg.Time.String())
		}

		msgKafka := &sarama.ConsumerMessage{
			Headers: []*sarama.RecordHeader{
				{[]byte(quark.HeaderMessageTime), publishTime},
				{[]byte(quark.HeaderMessageType), []byte("chat.0")},
			},
			Timestamp:      time.Time{},
			BlockTimestamp: time.Time{},
			Key:            nil,
			Value:          []byte("hello there"),
			Topic:          "",
			Partition:      0,
			Offset:         0,
		}

		msgMock := new(quark.Message)
		UnmarshalKafkaHeaders(msgKafka.Headers, msgMock)

		assert.Equal(t, msg.Time.String(), msgMock.Time.String())
		assert.Equal(t, msg.Type, msgMock.Type)
	})
}

func TestUnmarshalKafkaMessage(t *testing.T) {
	t.Run("Unmarshal Kafka message", func(t *testing.T) {
		msg := quark.NewMessage("1", "chat.0", []byte("hello there"))

		publishTime, err := msg.Time.MarshalBinary()
		if err != nil {
			publishTime = []byte(msg.Time.String())
		}

		msgKafka := &sarama.ConsumerMessage{
			Headers: []*sarama.RecordHeader{
				{[]byte(quark.HeaderMessageTime), publishTime},
				{[]byte(quark.HeaderMessageType), []byte("chat.0")},
				{[]byte(quark.HeaderMessageSpecVersion), []byte(quark.CloudEventsVersion)},
			},
			Timestamp:      time.Time{},
			BlockTimestamp: time.Time{},
			Key:            nil,
			Value:          []byte("hello there"),
			Topic:          "",
			Partition:      0,
			Offset:         0,
		}

		msgMock := new(quark.Message)
		UnmarshalKafkaMessage(msgKafka, msgMock)

		assert.Equal(t, msg.Time.String(), msgMock.Time.String())
		assert.Equal(t, msg.Type, msgMock.Type)
		assert.Equal(t, msg.SpecVersion, msgMock.SpecVersion)
	})
}
