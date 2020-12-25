package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	t.Run("Event Header item insertion", func(t *testing.T) {
		h := Header{}
		h.Set(HeaderMessageCorrelationId, "123")
		assert.Equal(t, "123", h.Get(HeaderMessageCorrelationId))
	})
	t.Run("Event Header verify if item exists", func(t *testing.T) {
		h := Header{}
		h.Set(HeaderMessageId, "1")
		assert.Equal(t, true, h.Contains(HeaderMessageId))
	})
	t.Run("Event Header item removal", func(t *testing.T) {
		h := Header{}
		h.Set(HeaderMessageId, "1")
		h.Del(HeaderMessageId)
		assert.Equal(t, false, h.Contains(HeaderMessageId))
	})
}

func BenchmarkHeader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Run("Event Header data insertion", func(b *testing.B) {
			b.ReportAllocs()
			h := Header{}
			h.Set(HeaderMessageId, "123")
		})
		b.Run("Event Header data removal", func(b *testing.B) {
			b.ReportAllocs()
			h := Header{}
			h.Set(HeaderMessageId, "123")
			h.Del(HeaderMessageId)
		})
		b.Run("Event Header contains item", func(b *testing.B) {
			b.ReportAllocs()
			h := Header{}
			h.Set(HeaderMessageId, "123")
			h.Contains(HeaderMessageId)
		})
		b.Run("Event Header data manipulation", func(b *testing.B) {
			b.ReportAllocs()
			h := Header{}
			h.Set(HeaderMessageId, "123")
			h.Get(HeaderMessageId)
			h.Contains(HeaderMessageId)
			h.Del(HeaderMessageId)
		})
	}
}
