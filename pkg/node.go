package pkg

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eapache/queue"
	"github.com/hashicorp/go-multierror"
)

// node is a Consumer logical unit of work.
//
// Distributes blocking I/O operations into different goroutines to enable parallelism with fan-out mechanisms.
type node struct {
	Broker   *Broker
	Consumer *Consumer

	workers        sync.Pool
	runningWorkers *queue.Queue
}

func newNode(b *Broker, c *Consumer) *node {
	n := &node{
		Broker:         b,
		Consumer:       c,
		workers:        sync.Pool{},
		runningWorkers: queue.New(),
	}
	n.workers = sync.Pool{
		New: func() interface{} {
			return newWorker(n) // returns a new worker when required
		},
	}

	return n
}

func (n *node) Consume(ctx context.Context) error {
	if err := n.ensureValidParams(); err != nil {
		return err
	}
	errs := new(multierror.Error)
	// Start worker jobs, these are Blocking I/O and each working should create a new goroutine.
	//
	// We cannot create goroutines from here because sync.Pool malloc reference is pooled back if we do so.
	// Besides, starting a worker job is already thread-safe.
	for i := 0; i < n.setDefaultPoolSize(); i++ {
		if w, ok := n.workers.Get().(worker); w != nil && ok {
			workerCtx := ctx
			if err := w.StartJob(workerCtx); err != nil {
				errs = multierror.Append(errs, err)
			} else {
				n.runningWorkers.Add(w)
			}
			continue
		}
		errs = multierror.Append(errs, fmt.Errorf("topic(s) %s: %w", n.Consumer.TopicString(),
			ErrProviderNotValid))
	}

	return errs.ErrorOrNil()
}

func (n *node) ensureValidParams() error {
	if len(n.setDefaultCluster()) == 0 {
		return ErrEmptyCluster
	} else if len(n.Consumer.topics) == 0 {
		return ErrNotEnoughTopics
	} else if n.Consumer.handlerFunc == nil && n.Consumer.handler == nil {
		return ErrNotEnoughHandlers
	}
	return nil
}

func (n *node) Close() error {
	errs := new(multierror.Error)
	runningLength := n.runningWorkers.Length() // allocate in a different memory address to avoid queue length mutation
	for i := 0; i < runningLength; i++ {
		w := n.runningWorkers.Remove().(worker)
		if err := w.Close(); err != nil {
			errs = multierror.Append(errs, err)
		}
		n.workers.Put(w) // avoid memory leaks by sending back workers to the pool
	}
	return errs.ErrorOrNil()
}

func (n *node) setDefaultPoolSize() int {
	if n.Consumer.poolSize > 0 {
		return n.Consumer.poolSize
	}
	return n.Broker.setDefaultPoolSize() // use global
}

func (n *node) setDefaultMaxRetries() int {
	if n.Consumer.maxRetries >= 0 {
		return n.Consumer.maxRetries
	}
	return n.Broker.setDefaultMaxRetries() // use global
}

func (n *node) setDefaultRetryBackoff() time.Duration {
	if n.Consumer.retryBackoff >= 0 {
		return n.Consumer.retryBackoff
	}
	return n.Broker.setDefaultRetryBackoff() // use global
}

func (n *node) setDefaultProvider() string {
	if provider := n.Consumer.provider; provider != "" {
		return provider
	}
	return n.Broker.Provider // use global
}

func (n *node) setDefaultProviderConfig() interface{} {
	if cfg := n.Consumer.providerConfig; cfg != nil {
		return cfg
	}
	return n.Broker.ProviderConfig // use global
}

func (n *node) setDefaultCluster() []string {
	if cluster := n.Consumer.cluster; len(cluster) > 0 {
		return cluster
	}
	return n.Broker.Cluster // use global
}

func (n *node) setDefaultGroup() string {
	if group := n.Consumer.group; group != "" {
		return group
	}
	return n.Consumer.TopicString() // use topics as default
}

func (n *node) setDefaultPublisher() Publisher {
	if n.Consumer.publisher != nil {
		return n.Consumer.publisher
	} else if n.Broker.Publisher != nil {
		return n.Broker.Publisher // use global
	}

	return getDefaultPublisher(n.setDefaultProvider())
}

func (n *node) setDefaultEventWriter() EventWriter {
	if n.Broker.EventWriter != nil {
		return n.Broker.EventWriter
	}

	return newEventWriter(n, n.setDefaultPublisher())
}
