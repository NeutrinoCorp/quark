package pkg

import (
	"context"

	"github.com/Shopify/sarama"
)

type worker interface {
	StartJob(context.Context)
}

func newWorker(n *node, provider string, providerCfg interface{}) worker {
	if provider == "" {
		return nil
	}
	// strategy pattern
	switch provider {
	case KafkaProvider:
		if cfg, ok := providerCfg.(KafkaConfiguration); ok {
			return kafkaWorker{
				node: n,
				cfg:  cfg.Config,
			}
		}

	case AWSProvider:
		return nil
	}

	return nil
}

type kafkaWorker struct {
	node *node
	cfg  *sarama.Config
}

func (k kafkaWorker) StartJob(ctx context.Context) {
	panic("implement me")
}
