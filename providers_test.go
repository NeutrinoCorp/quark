package quark

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var providerTestingSuite = []struct {
	provider string
	cfg      interface{}
	expected error
}{
	{KafkaProvider, KafkaConfiguration{}, nil},
	{AWSProvider, AWSConfiguration{}, nil},
	{KafkaProvider, nil, ErrProviderNotValid},
	{AWSProvider, nil, ErrProviderNotValid},
	{KafkaProvider, AWSConfiguration{}, ErrProviderNotValid},
	{AWSProvider, KafkaConfiguration{}, ErrProviderNotValid},
	{ActiveMqProvider, nil, ErrProviderNotValid},
	{RabbitMqProvider, nil, ErrProviderNotValid},
	{NatsProvider, nil, ErrProviderNotValid},
	{"", nil, ErrProviderNotValid},
}

func TestProviders(t *testing.T) {
	for _, tt := range providerTestingSuite {
		t.Run("Ensure valid provider", func(t *testing.T) {
			err := ensureValidProvider(tt.provider, tt.cfg)
			assert.True(t, errors.Is(err, tt.expected))
		})
	}
}

func BenchmarkProviders(b *testing.B) {
	cfg := AWSConfiguration{}
	b.Run("Ensure valid provider", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ensureValidProvider(AWSProvider, cfg)
		}
	})
}
