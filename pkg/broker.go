package pkg

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/Shopify/sarama"
)

// Broker coordinates every Event operation on the running Event-Driven application.
//
// Administrates Consumer(s) and its worker nodes wrapped with well-known concurrency and resiliency patterns.
type Broker struct {
	Provider       string
	ProviderConfig interface{}
	Cluster        []string
	ErrorHandler   func(error)
	Publisher      Publisher
	EventMux       EventMux
	PoolSize       int
	MaxRetries     int
	RetryBackoff   time.Duration

	nodes sync.Pool
}

var defaultPoolSize = 5

// NewBroker allocates and returns a Broker
func NewBroker() *Broker {
	return &Broker{
		Provider:       "",
		ProviderConfig: nil,
		Cluster:        make([]string, 0),
		ErrorHandler:   nil,
		Publisher:      nil,
		EventMux:       nil,
	}
}

// NewKafkaBroker allocates and returns a Kafka Broker
func NewKafkaBroker(cfg *sarama.Config, addrs ...string) *Broker {
	return &Broker{
		Provider: "kafka",
		ProviderConfig: KafkaConfiguration{
			Config: cfg,
		},
		Cluster:      addrs,
		ErrorHandler: nil,
		Publisher:    nil,
		EventMux:     nil,
	}
}

func (b *Broker) setDefaultMux() {
	if b.EventMux == nil {
		b.EventMux = NewMux()
	}
}

// ListenAndServe starts listening to the given Consumer(s) concurrently-safe
func (b *Broker) ListenAndServe() error {
	b.setDefaultMux()
	ctx := context.Background()
	errs := new(multierror.Error)
	for _, c := range b.EventMux.List() {
		n := &node{
			Broker:   b,
			Consumer: c,
			workers:  sync.Pool{},
		}
		if err := n.Populate(); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		n.Consume(ctx)
		b.nodes.Put(n)
	}
	return errs.ErrorOrNil()
}

// Topic adds new Consumer node to the given EventMux
func (b *Broker) Topic(topic string) *Consumer {
	b.setDefaultMux()
	return b.EventMux.Topic(topic)
}

// Topics adds multiple Consumer nodes to the given EventMux
func (b *Broker) Topics(topics ...string) *Consumer {
	b.setDefaultMux()
	return b.EventMux.Topics(topics...)
}

func (b *Broker) setDefaultPoolSize() int {
	if b.PoolSize > 0 {
		return b.PoolSize
	}
	return defaultPoolSize
}
