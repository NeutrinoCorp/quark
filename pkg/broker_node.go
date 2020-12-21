package pkg

import (
	"context"
	"sync"
)

type node struct {
	Broker   *Broker
	Consumer *Consumer

	workers sync.Pool
}

func (n *node) Populate() error {
	w := n.setDefaultWorker()
	if w == nil {
		return ErrProviderNotValid
	}

	for i := 0; i < n.setDefaultPoolSize(); i++ {
		n.workers.Put(w)
	}
	return nil
}

func (n *node) Consume(ctx context.Context) {
	for i := 0; i < n.setDefaultPoolSize(); i++ {
		if w, ok := n.workers.Get().(worker); ok && w != nil {
			go w.StartJob(ctx)
		}
	}
}

func (n *node) setDefaultWorker() worker {
	var w worker
	if w = newWorker(n, n.Consumer.provider, n.Consumer.Config()); w == nil {
		w = newWorker(n, n.Broker.Provider, n.Broker.ProviderConfig) // use global provider
	}
	return w
}

func (n *node) setDefaultPoolSize() int {
	if n.Consumer.poolSize > 0 {
		return n.Consumer.poolSize
	}
	return n.Broker.setDefaultPoolSize() // use global
}
