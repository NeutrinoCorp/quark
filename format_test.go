package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatTopicName(t *testing.T) {
	t.Run("Comply with Topic Async API specification", func(t *testing.T) {
		correctName := "neutrino.payment.1.domain_event.user.paid"
		topic := FormatTopicName("neutrino", "payment", DomainEvent, "user", "paid", 1)
		assert.Equal(t, correctName, topic)
	})
}

func BenchmarkFormatTopicName(b *testing.B) {
	b.Run("Form a valid Async API topic name", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatTopicName("neutrino", "payment", DomainEvent, "user", "paid", 1)
		}
	})
}

func TestFormatQueueName(t *testing.T) {
	t.Run("Comply with Queue Async API specification", func(t *testing.T) {
		correctName := "payment.organization.increment_sales_on_user_paid"
		topic := FormatQueueName("payment", "organization", "increment_sales", "user_paid")
		assert.Equal(t, correctName, topic)
	})
}

func BenchmarkFormatQueueName(b *testing.B) {
	b.Run("Form a valid Async API queue name", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatQueueName("payment", "organization", "increment_sales", "user_paid")
		}
	})
}
