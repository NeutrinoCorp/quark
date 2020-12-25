package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var consumerAddTopicsTestingSuite = []struct {
	n        []string // input
	expected int      // expected result (length)
}{
	{[]string{"chat.0"}, 1},
	{[]string{"chat.0", "chat.1"}, 2},
	{[]string{}, 0},
	{[]string{"chat.0", "chat.1", "chat.2"}, 3},
}

func TestConsumer(t *testing.T) {
	t.Run("Consumer add topics", func(t *testing.T) {
		for _, tt := range consumerAddTopicsTestingSuite {
			c := Consumer{}
			c.Topics(tt.n...)
			assert.Equal(t, tt.expected, len(c.topics))
		}
	})
	t.Run("Consumer public operations", func(t *testing.T) {
		c := Consumer{}
		c.Topic("chat.1")
		assert.Contains(t, c.topics, "chat.1")
		c.Provider(KafkaProvider)
		assert.Equal(t, "kafka", c.provider)
		c.PoolSize(10)
		assert.Equal(t, 10, c.poolSize)
		c.ProviderConfig(KafkaConfiguration{})
		assert.Equal(t, KafkaConfiguration{}, c.providerConfig.(KafkaConfiguration))
		c.Address("localhost")
		assert.Contains(t, c.cluster, "localhost")
	})
}

func BenchmarkConsumer(b *testing.B) {
	topics := []string{"chat.0", "chat.1", "chat.2", "chat.3"}
	b.Run("Consumer add topics", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			c := Consumer{}
			c.Topics(topics...)
		}
	})
	b.Run("Consumer public operations", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			c := Consumer{}
			c.Topic("chat.1")
			c.Provider(KafkaProvider)
			c.PoolSize(10)
			c.ProviderConfig(KafkaConfiguration{})
			c.Address("localhost")
		}
	})
}
