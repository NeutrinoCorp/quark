package quark

import (
	"strconv"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func TestNewKafkaHeader(t *testing.T) {
	t.Run("New kafka header", func(t *testing.T) {
		now := time.Now()
		kafkaMsg := &sarama.ConsumerMessage{}
		kafkaMsg.Partition = 1
		kafkaMsg.Offset = 30
		kafkaMsg.Key = []byte("kafka")
		kafkaMsg.Value = []byte("message content")
		kafkaMsg.Timestamp = now
		kafkaMsg.BlockTimestamp = now.Add(time.Hour * 1)
		kafkaMsg.Headers = []*sarama.RecordHeader{
			{[]byte(HeaderMessageId), []byte("1")},
		}

		h := NewKafkaHeader(kafkaMsg)
		assert.Equal(t, strconv.Itoa(int(kafkaMsg.Partition)), h.Get(HeaderKafkaPartition))
		assert.Equal(t, strconv.Itoa(int(kafkaMsg.Offset)), h.Get(HeaderKafkaOffset))
		assert.Equal(t, string(kafkaMsg.Key), h.Get(HeaderKafkaKey))
		assert.Equal(t, string(kafkaMsg.Value), h.Get(HeaderKafkaValue))
		assert.Equal(t, kafkaMsg.Timestamp.String(), h.Get(HeaderKafkaTimestamp))
		assert.Equal(t, kafkaMsg.BlockTimestamp.String(), h.Get(HeaderKafkaBlockTimestamp))
		assert.Equal(t, string(kafkaMsg.Headers[0].Value), h.Get(HeaderMessageId))
	})
}

func BenchmarkNewKafkaHeader(b *testing.B) {
	now := time.Now()
	kafkaMsg := &sarama.ConsumerMessage{}
	kafkaMsg.Partition = 1
	kafkaMsg.Offset = 30
	kafkaMsg.Key = []byte("kafka")
	kafkaMsg.Value = []byte("message content")
	kafkaMsg.Timestamp = now
	kafkaMsg.BlockTimestamp = now.Add(time.Hour * 1)
	kafkaMsg.Headers = []*sarama.RecordHeader{
		{[]byte(HeaderMessageId), []byte("1")},
	}
	b.Run("New kafka header", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = NewKafkaHeader(kafkaMsg)
		}
	})
}
