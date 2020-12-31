package quark

import "sync"

// EventMux Consumer node registry, the mux itself contains a local registry of all the given Consumer(s).
// This process is required by the Broker in order to start all the nodes registered on the mux.
type EventMux interface {
	// Topic adds a new topic to the local registry
	Topic(topic string) *Consumer
	// Topics adds multiple topics
	Topics(topics ...string) *Consumer
	// Add manually adds a consumer
	Add(c *Consumer)
	// Get returns an specific consumer by the given topic
	Get(topic string) []*Consumer
	// Del removes an specific topic from the local registry
	Del(topic string)
	// Contains verifies if the given topic exists in the local registry
	Contains(topic string) bool
	// List returns the local registry
	List() map[string][]*Consumer
}

type defaultMux struct {
	consumers map[string][]*Consumer
	mu        sync.RWMutex
}

// NewMux allocates and returns a default EventMux
func NewMux() EventMux {
	return &defaultMux{
		consumers: map[string][]*Consumer{},
		mu:        sync.RWMutex{},
	}
}

func (b *defaultMux) Topic(topic string) *Consumer {
	b.mu.Lock()
	defer b.mu.Unlock()
	c := new(Consumer)
	if topic == "" {
		return c
	}
	c.Topic(topic)
	b.consumers[topic] = append(b.consumers[topic], c)
	return c
}

func (b *defaultMux) Topics(topics ...string) *Consumer {
	b.mu.Lock()
	defer b.mu.Unlock()
	c := new(Consumer)
	if len(topics) == 0 {
		return c
	}
	c.Topics(topics...)
	// b.consumers.Store(c.TopicString(), c) -  previous version, attached many topics into a single unit of work
	// Adds a worker pool per-topic
	for _, t := range topics {
		b.consumers[t] = append(b.consumers[t], c)
	}
	return c
}

func (b *defaultMux) Add(c *Consumer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if c == nil {
		return // ignore nil-refs
	}
	for _, t := range c.topics {
		b.consumers[t] = append(b.consumers[t], c)
	}
}

func (b *defaultMux) Get(topic string) []*Consumer {
	b.mu.RLock()
	defer b.mu.RUnlock()
	c, ok := b.consumers[topic]
	if !ok {
		return nil
	}

	return c
}

func (b *defaultMux) Del(topic string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.consumers, topic)
}

func (b *defaultMux) Contains(topic string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.consumers[topic]
	return ok
}

func (b *defaultMux) List() map[string][]*Consumer {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.consumers
}
