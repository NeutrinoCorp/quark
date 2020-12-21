package pkg

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
	Get(topic string) *Consumer
	// Del removes an specific topic from the local registry
	Del(topic string)
	// Contains verifies if the given topic exists in the local registry
	Contains(topic string) bool
	// List returns the local registry
	List() map[string]*Consumer
}

type defaultMux struct {
	consumers sync.Map
}

// NewMux allocates and returns a default EventMux
func NewMux() EventMux {
	return new(defaultMux)
}

func (b *defaultMux) Topic(topic string) *Consumer {
	c := new(Consumer)
	c.Topic(topic)
	b.consumers.Store(topic, c)
	return c
}

func (b *defaultMux) Topics(topics ...string) *Consumer {
	c := new(Consumer)
	c.Topics(topics...)
	for _, t := range topics {
		b.consumers.Store(t, c)
	}
	return c
}

func (b *defaultMux) Add(c *Consumer) {
	for _, t := range c.topics {
		b.consumers.Store(t, c)
	}
}

func (b *defaultMux) Get(topic string) *Consumer {
	cSync, ok := b.consumers.Load(topic)
	if c, okC := cSync.(*Consumer); ok && okC {
		return c
	}

	return nil
}

func (b *defaultMux) Del(topic string) {
	b.consumers.Delete(topic)
}

func (b *defaultMux) Contains(topic string) bool {
	_, ok := b.consumers.Load(topic)
	return ok
}

func (b *defaultMux) List() map[string]*Consumer {
	list := map[string]*Consumer{}
	b.consumers.Range(func(key, value interface{}) bool {
		k, okK := key.(string)
		v, okV := value.(*Consumer)
		if okK && okV {
			list[k] = v
		}
		return true
	})
	return list
}
