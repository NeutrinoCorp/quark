package quark

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
)

// Broker coordinates every Event operation on the running Event-Driven application.
//
// Administrates Consumer(s) Nodes and their workers wrapped with well-known concurrency and resiliency patterns.
type Broker struct {
	Provider         string
	ProviderConfig   interface{}
	Cluster          []string
	ErrorHandler     func(context.Context, error)
	Publisher        Publisher
	EventMux         EventMux
	EventWriter      EventWriter
	PoolSize         int
	MaxRetries       int
	RetryBackoff     time.Duration
	ConnRetries      int
	ConnRetryBackoff time.Duration

	MessageIdGenerator func() string
	WorkerFactory      WorkerFactory

	BaseContext context.Context

	nodes          map[int]*Node
	runningNodes   int
	runningWorkers int
	mu             sync.Mutex
	inShutdown     atomicBool
	doneChan       chan struct{}
}

var (
	defaultPoolSize         = 5
	defaultMaxRetries       = 1
	defaultRetryBackoff     = time.Second * 3
	defaultConnRetries      = 3
	defaultConnRetryBackoff = time.Second * 5
	shutdownPollInterval    = time.Millisecond * 500
)

// NewBroker allocates and returns a Broker
func NewBroker(provider string, config interface{}) *Broker {
	return &Broker{
		Provider:       provider,
		ProviderConfig: config,
		Cluster:        make([]string, 0),
		ErrorHandler:   nil,
		Publisher:      nil,
		EventMux:       nil,
		EventWriter:    nil,
		nodes:          make(map[int]*Node),
		mu:             sync.Mutex{},
		inShutdown:     0,
		doneChan:       nil,
	}
}

// ListenAndServe starts listening to the given Consumer(s) concurrently-safe
func (b *Broker) ListenAndServe() error {
	if b.shuttingDown() {
		return ErrBrokerClosed
	}
	return b.Serve()
}

// Serve starts the broker components
func (b *Broker) Serve() error {
	for {
		if b.BaseContext == nil {
			b.BaseContext = context.Background()
		}
		b.setDefaultMux()
		if err := b.startNodes(b.BaseContext); err != nil {
			return err
		}

		<-b.getDoneChanLocked()
		b.Shutdown(b.BaseContext)
	}
}

func (b *Broker) startNodes(ctx context.Context) error {
	for _, consumers := range b.EventMux.List() {
		for _, c := range consumers {
			nodeCtx := ctx
			n := newNode(b, c)
			if err := n.Consume(nodeCtx); err != nil {
				return err
			}
			b.nodes[b.runningNodes] = n
			b.runningWorkers += n.runningWorkers.Length()
			b.runningNodes++
		}
	}
	return nil
}

// Shutdown starts Broker graceful shutdown of its components
func (b *Broker) Shutdown(ctx context.Context) error {
	b.inShutdown.setTrue()
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closeDoneChanLocked()

	ticker := time.NewTicker(shutdownPollInterval)
	defer ticker.Stop()
	for {
		if err := b.closeNodes(); err == nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (b *Broker) shuttingDown() bool {
	return b.inShutdown.isSet()
}

func (b *Broker) getDoneChanLocked() chan struct{} {
	if b.doneChan == nil {
		b.doneChan = make(chan struct{})
	}
	return b.doneChan
}

func (b *Broker) closeDoneChanLocked() {
	ch := b.getDoneChanLocked()
	select {
	case <-ch:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.
		close(ch)
	}
}

func (b *Broker) closeNodes() error {
	errs := new(multierror.Error)
	for k, n := range b.nodes {
		if err := n.Close(); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		b.runningNodes--
		b.runningWorkers = n.runningWorkers.Length()
		delete(b.nodes, k)
	}

	return errs.ErrorOrNil()
}

// Topic adds new Consumer Node to the given EventMux
func (b *Broker) Topic(topic string) *Consumer {
	b.setDefaultMux()
	return b.EventMux.Topic(topic)
}

// Topics adds multiple Consumer nodes to the given EventMux
func (b *Broker) Topics(topics ...string) *Consumer {
	b.setDefaultMux()
	return b.EventMux.Topics(topics...)
}

// RunningNodes returns the current number of running nodes
func (b *Broker) RunningNodes() int {
	return b.runningNodes
}

// RunningWorkers returns the current number of running workers (inside every Node)
func (b *Broker) RunningWorkers() int {
	return b.runningWorkers
}

func (b *Broker) setDefaultMux() {
	if b.EventMux == nil {
		b.EventMux = NewMux()
	}
}

func (b *Broker) setDefaultPoolSize() int {
	if b.PoolSize > 0 {
		return b.PoolSize
	}
	return defaultPoolSize
}

func (b *Broker) setDefaultMaxRetries() int {
	if b.MaxRetries > 0 {
		return b.MaxRetries
	}
	return defaultMaxRetries
}

func (b *Broker) setDefaultRetryBackoff() time.Duration {
	if b.RetryBackoff > 0 {
		return b.RetryBackoff
	}
	return defaultRetryBackoff
}

func (b *Broker) setDefaultConnRetries() int {
	if b.ConnRetries > 0 {
		return b.ConnRetries
	}
	return defaultConnRetries
}

func (b *Broker) setDefaultConnRetryBackoff() time.Duration {
	if b.ConnRetryBackoff > 0 {
		return b.ConnRetryBackoff
	}
	return defaultConnRetryBackoff
}

func (b *Broker) setDefaultMessageIdGenerator() IdGenerator {
	if b.MessageIdGenerator != nil {
		return b.MessageIdGenerator
	}
	return defaultIdGenerator
}

// GetConnRetries retrieves the default connection retries
func (b *Broker) GetConnRetries() int {
	return b.setDefaultConnRetries()
}

// GetConnRetryBackoff retrieves the default connection retry backoff
func (b *Broker) GetConnRetryBackoff() time.Duration {
	return b.setDefaultConnRetryBackoff()
}
