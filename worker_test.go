package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var workerTestingSuite = []struct {
	n        string // input
	cfg      interface{}
	expected string // expected result
}{
	{"", nil, ""},
	{"KaFkA", nil, ""},
	{"Apache Kafka", KafkaConfiguration{}, ""},
	{"kafka", AWSConfiguration{}, ""},
	{KafkaProvider, KafkaConfiguration{}, KafkaProvider},
	{"AWS", nil, ""},
	{"Amazon Web Services", AWSConfiguration{}, ""},
	{"aws", KafkaConfiguration{}, ""},
	{AWSProvider, AWSConfiguration{}, AWSProvider},
}

func TestWorker(t *testing.T) {
	for _, tt := range workerTestingSuite {
		t.Run("New worker", func(t *testing.T) {
			b, err := NewBroker(KafkaProvider, KafkaConfiguration{})
			assert.Nil(t, err)
			b.Topic("foo.1")
			var n *node
			for _, consumers := range b.EventMux.List() {
				for _, c := range consumers {
					n = newNode(b, c)
				}
			}
			assert.NotNil(t, n)
			n.Consumer.Provider(tt.n)
			n.Consumer.ProviderConfig(tt.cfg)
			// bypass broker Insurance
			n.Broker.Provider = tt.n
			n.Broker.ProviderConfig = tt.cfg
			w := newWorker(n)
			if w != nil {
				assert.Equal(t, tt.expected, w.Parent().Consumer.provider)
				return
			}
			assert.Empty(t, w)
		})
	}
}

func BenchmarkWorker(b *testing.B) {
	b.Run("New worker", func(b *testing.B) {
		b.ReportAllocs()
		br, err := NewBroker(KafkaProvider, KafkaConfiguration{})
		if err != nil {
			b.Error(err)
		}
		br.Topic("foo.1")
		var n *node
		for _, consumers := range br.EventMux.List() {
			for _, c := range consumers {
				n = newNode(br, c)
			}
		}
		for i := 0; i < b.N; i++ {
			w := newWorker(n)
			_ = w.Parent().Broker.Provider
		}
	})
}
