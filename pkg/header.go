package pkg

// Header holds relevant metadata about the Message, Consumer or Broker.
//
//
// For example, a header may contain the following relevant information:
//
// - Broker and Consumer configurations.
//
// - Specific provider data (e.g. offsets, partitions and member id from a consumer group in Kafka).
//
// - An span context from the distributed tracing provider, so it can propagate traces from remote parents.
//
// - Message properties (e.g. total redeliveries, correlation id, origin service ip address).
type Header map[string]string

// Set attach or override the given fields
func (h Header) Set(k, v string) {
	h[k] = v
}

// Get returns a value from the given key
func (h Header) Get(k string) string {
	return h[k]
}

// Del removes a record
func (h Header) Del(k string) {
	delete(h, k)
}

// Contains verifies if the given record exists
func (h Header) Contains(k string) bool {
	_, ok := h[k]
	return ok
}
