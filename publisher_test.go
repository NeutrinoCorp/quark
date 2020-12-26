package quark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var publisherTestingSuite = []struct {
	provider string
	cfg      interface{}
	cluster  []string
	expected interface{}
}{
	{KafkaProvider, KafkaConfiguration{}, []string{"localhost:9092"}, &defaultKafkaPublisher{}},
	{AWSProvider, AWSConfiguration{}, []string{"https://sns.aws.com/my-aws-sns-endpoint"}, nil},
	{KafkaProvider, nil, []string{}, nil},
	{AWSProvider, nil, []string{}, nil},
	{KafkaProvider, AWSConfiguration{}, []string{}, nil},
	{AWSProvider, KafkaConfiguration{}, []string{}, nil},
	{ActiveMqProvider, nil, []string{}, nil},
	{RabbitMqProvider, nil, []string{}, nil},
	{NatsProvider, nil, []string{}, nil},
	{"", nil, []string{}, nil},
}

func TestDefaultPublisher(t *testing.T) {
	for _, tt := range publisherTestingSuite {
		t.Run("Default publisher", func(t *testing.T) {
			p := getDefaultPublisher(tt.provider, tt.cfg, tt.cluster)
			if tt.expected != nil {
				assert.NotNil(t, p)
				return
			}
			assert.Nil(t, p)
		})
	}
}

func BenchmarkDefaultPublisher(b *testing.B) {
	cfg := KafkaConfiguration{}
	cluster := []string{"localhost:9092"}
	b.Run("Default publisher", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = getDefaultPublisher(KafkaProvider, cfg, cluster)
		}
	})
}
