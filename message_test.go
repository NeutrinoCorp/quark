package quark

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var newMessageTestingSuite = []struct {
	n        string // input
	expected string // expected result
}{
	{"chat.0", "chat.0"},
	{"alex.trades", "alex.trades"},
	{"alice.notifications", "alice.notifications"},
	{"bob.gps", "bob.gps"},
}

func TestNewMessage(t *testing.T) {
	id := uuid.New().String()
	t.Run("New Message", func(t *testing.T) {
		for _, tt := range newMessageTestingSuite {
			msg := NewMessage(id, tt.n, []byte("hello"))
			assert.Equal(t, tt.expected, msg.Kind)
			assert.Equal(t, "hello", string(msg.Attributes))
		}
	})
}

func BenchmarkNewMessage(b *testing.B) {
	id := uuid.New().String()
	b.Run("New message", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = NewMessage(id, "chat.0", []byte("hello"))
		}
	})
}

func TestNewMessageFromParent(t *testing.T) {
	parentId := uuid.New().String()
	id := uuid.New().String()
	t.Run("New Message from parent", func(t *testing.T) {
		for _, tt := range newMessageTestingSuite {
			msg := NewMessageFromParent(parentId, id, tt.n, []byte("hello"))
			assert.Equal(t, tt.expected, msg.Kind)
			assert.Equal(t, "hello", string(msg.Attributes))
			assert.Equal(t, parentId, msg.Metadata.CorrelationId)
		}
	})
}

func BenchmarkNewMessageFromParent(b *testing.B) {
	parentId := uuid.New().String()
	id := uuid.New().String()
	b.Run("New message from parent", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = NewMessageFromParent(parentId, id, "chat.0", []byte("hello"))
		}
	})
}
