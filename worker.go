package quark

import (
	"context"
)

// worker is a unit of work a Consumer node.
//
// It handles all the work and should be running inside a goroutine to enable parallelism.
type worker interface {
	// SetID inject the current worker id
	SetID(i int)
	// Parent returns the parent node
	Parent() *node
	// StartJob starts an specific work
	//
	// Panics goroutine if any errors is thrown
	StartJob(context.Context) error
	// Close stop all Blocking I/O operations
	Close() error
}

func newWorker(parent *node) worker {
	provider := parent.setDefaultProvider()
	if provider == "" {
		return nil // avoid extra computation by this return
	}
	providerCfg := parent.setDefaultProviderConfig()
	if providerCfg == nil {
		return nil // avoid extra computation by this return
	}

	// strategy
	switch provider {
	case KafkaProvider:
		if cfg, ok := providerCfg.(KafkaConfiguration); ok {
			return &kafkaWorker{
				id:     0,
				parent: parent,
				cfg:    cfg,
			}
		}
		return nil
	case AWSProvider:
		if cfg, ok := providerCfg.(AWSConfiguration); ok {
			return &awsWorker{
				id:     0,
				parent: parent,
				cfg:    cfg,
			}
		}
		return nil
	default:
		return nil
	}
}
