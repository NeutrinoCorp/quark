package quark

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type stubHandler struct{}

func (h stubHandler) ServeEvent(EventWriter, *Event) bool {
	return true
}

var consumerAddTopicsTestingSuite = []struct {
	n        []string // input
	expected int      // expected result (length)
}{
	{[]string{"chat.0"}, 1},
	{[]string{"chat.0", "chat.1"}, 2},
	{[]string{}, 0},
	{[]string{"chat.0", "chat.1", "chat.2"}, 3},
}

var consumerAddTopicTestingSuite = []struct {
	t        string // input
	expected int    // expected result (length)
}{
	{"chat.0", 1},
	{"", 0},
}

func TestConsumer(t *testing.T) {
	t.Run("Consumer add topics", func(t *testing.T) {
		for _, tt := range consumerAddTopicsTestingSuite {
			c := Consumer{}
			c.Topics(tt.n...)
			assert.Equal(t, tt.expected, len(c.topics))
		}
	})
	t.Run("Consumer add topic", func(t *testing.T) {
		for _, tt := range consumerAddTopicTestingSuite {
			c := Consumer{}
			c.Topic(tt.t)
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
		c.Address("localhost")
		assert.Contains(t, c.cluster, "localhost")
	})
}

func TestConsumer_Group(t *testing.T) {
	t.Run("Consumer group mutation", func(t *testing.T) {
		c := Consumer{}
		c.Group("foo-group")
		assert.Equal(t, "foo-group", c.group)
	})
}

func TestConsumer_MaxRetries(t *testing.T) {
	t.Run("Consumer max retries mutation", func(t *testing.T) {
		c := Consumer{}
		c.MaxRetries(5)
		assert.Equal(t, 5, c.maxRetries)
	})
}

func TestConsumer_RetryBackoff(t *testing.T) {
	t.Run("Consumer retry backoff mutation", func(t *testing.T) {
		c := Consumer{}
		c.RetryBackoff(time.Second * 5)
		assert.Equal(t, time.Second*5, c.retryBackoff)
	})
}

type stubProviderConfig struct {
	Foo string
}

func TestConsumer_GetProviderConfig(t *testing.T) {
	t.Run("Consumer provider config mutation", func(t *testing.T) {
		c := Consumer{}
		cfg := &stubProviderConfig{}
		c.ProviderConfig(cfg)
		assert.Same(t, cfg, c.GetProviderConfig())
	})
}

func TestConsumer_Handle(t *testing.T) {
	t.Run("Consumer handle mutation", func(t *testing.T) {
		c := Consumer{}
		h := &stubHandler{}
		c.Handle(h)
		assert.Same(t, h, c.handler)
	})
}

func TestConsumer_HandleFunc(t *testing.T) {
	t.Run("Consumer handle function mutation", func(t *testing.T) {
		c := Consumer{}
		h := func(EventWriter, *Event) bool { return true }
		c.HandleFunc(h)
		assert.NotEmpty(t, c.handlerFunc)
	})
}

var consumerTopicStringTestingSuite = []struct {
	topics   []string
	expected string
}{
	{[]string{"chat.0", "chat.1", "chat.2"}, "chat.0,chat.1,chat.2"},
	{[]string{"chat.0"}, "chat.0"},
	{[]string{}, ""},
}

func TestConsumer_TopicString(t *testing.T) {
	for _, tt := range consumerTopicStringTestingSuite {
		t.Run("Consumer topics string", func(t *testing.T) {
			c := Consumer{}
			c.Topics(tt.topics...)
			assert.Equal(t, tt.expected, c.TopicString())
		})
	}
}

type stubPublisher struct {
	fail bool
}

var errStubPublisher = errors.New("generic stub publisher error")

func (a stubPublisher) Publish(context.Context, ...*Message) error {
	if a.fail {
		return errStubPublisher
	}
	return nil
}

func TestConsumer_Publisher(t *testing.T) {
	t.Run("Consumer add publisher", func(t *testing.T) {
		c := Consumer{}
		p := new(stubPublisher)
		c.Publisher(p)
		assert.Equal(t, p, c.publisher)
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
			c.Address("localhost")
		}
	})
}
