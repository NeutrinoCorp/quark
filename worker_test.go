package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var workerTestingSuite = []struct {
	n        string // input
	expected string // expected result
}{
	{"KaFkA", ""},
	{"Apache Kafka", ""},
	{"kafka", "kafka"},
	{KafkaProvider, KafkaProvider},
	{"AWS", ""},
	{"Amazon Web Services", ""},
	{"aws", "aws"},
	{AWSProvider, AWSProvider},
}

func TestWorker(t *testing.T) {
	b, err := NewBroker(KafkaProvider, KafkaConfiguration{})
	if err != nil {
		t.Fatal(err.Error())
	}
	b.Topic("foo.1")
	var n *node
	for _, c := range b.EventMux.List() {
		n = newNode(b, c)
	}
	for _, tt := range workerTestingSuite {
		t.Run("New worker "+tt.n, func(t *testing.T) {
			n.Consumer.provider = tt.n
			w := newWorker(n)
			if w != nil {
				assert.Equal(t, tt.expected, w.Parent().Consumer.provider)
			}
		})
	}
}

func BenchmarkWorker(b *testing.B) {
	br, err := NewBroker(KafkaProvider, KafkaConfiguration{})
	if err != nil {
		b.Error(err)
	}
	br.Topic("foo.1")
	var n *node
	for _, c := range br.EventMux.List() {
		n = newNode(br, c)
	}
	for i := 0; i < b.N; i++ {
		b.Run("New worker", func(b *testing.B) {
			b.ReportAllocs()
			w := newWorker(n)
			_ = w.Parent().Broker.Provider
		})
	}
}
