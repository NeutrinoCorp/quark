package kafka

import (
	"strconv"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/neutrinocorp/quark"
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
			{[]byte(quark.HeaderMessageId), []byte("1")},
		}

		h := NewKafkaHeader(kafkaMsg)
		assert.Equal(t, strconv.Itoa(int(kafkaMsg.Partition)), h.Get(HeaderPartition))
		assert.Equal(t, strconv.Itoa(int(kafkaMsg.Offset)), h.Get(HeaderOffset))
		assert.Equal(t, string(kafkaMsg.Key), h.Get(HeaderKey))
		assert.Equal(t, string(kafkaMsg.Value), h.Get(HeaderValue))
		assert.Equal(t, kafkaMsg.Timestamp.String(), h.Get(HeaderTimestamp))
		assert.Equal(t, kafkaMsg.BlockTimestamp.String(), h.Get(HeaderBlockTimestamp))
		assert.Equal(t, string(kafkaMsg.Headers[0].Value), h.Get(quark.HeaderMessageId))
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
		{[]byte(quark.HeaderMessageId), []byte("1")},
	}
	b.Run("New kafka header", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = NewKafkaHeader(kafkaMsg)
		}
	})
}
