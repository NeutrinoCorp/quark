package quark

import "time"

// Consumer main processing unit. It is intended to subscribe to a Topic or Queue within a worker pool to
// segregate application load and enable high-concurrency.
//
// A Broker will receive and coordinate all Consumer nodes and will stop them gracefully when desired.
type Consumer struct {
	// Topic A Broker will use these topic(s) to subscribe the Consumer
	//	These field can be also used as Queues
	topics []string
	// Provider Message Broker or Messaging Queue system that this specific consumer will use (e.g. Kafka, RabbitMQ)
	provider string
	// ProviderConfig Custom provider configuration (e.g. sarama config, aws credentials)
	providerConfig interface{}
	// Group set of consumers this specific consumer must be with to consume messages in parallel
	//
	// 	Only available in: Apache Kafka
	group string
	// Cluster ip address(es) with its respective port(s) of the Message Broker/Message Queue system cluster
	cluster []string
	// Publisher pushes the given Message into the Event-Driven ecosystem.
	publisher Publisher
	// PoolSize worker pool size
	poolSize int
	// MaxRetries total times to retry consuming messages if processing fails
	maxRetries int
	// RetryBackoff time to wait between each retry
	retryBackoff time.Duration
	// Handler specific struct Quark will use to send messages
	handler Handler
	// HandlerFunc specific func Quark will use to send messages
	handlerFunc HandlerFunc
}

// Topic A Broker will use this topic to subscribe the Consumer
//	This field can be also used as Queue
func (c *Consumer) Topic(topic string) *Consumer {
	c.topics = append(c.topics, topic)
	return c
}

// Topics A Broker will use these topics to subscribe the Consumer to fan-in processing
//	These fields can be also used as Queues
func (c *Consumer) Topics(topics ...string) *Consumer {
	c.topics = append(c.topics, topics...)
	return c
}

// PoolSize worker pool size
func (c *Consumer) PoolSize(s int) *Consumer {
	c.poolSize = s
	return c
}

// MaxRetries total times to retry consuming messages if processing fails
func (c *Consumer) MaxRetries(n int) *Consumer {
	c.maxRetries = n
	return c
}

// RetryBackoff time to wait between each retry
func (c *Consumer) RetryBackoff(t time.Duration) *Consumer {
	c.retryBackoff = t
	return c
}

// Provider Message Broker or Messaging Queue system that this specific consumer will use (e.g. Kafka, RabbitMQ)
func (c *Consumer) Provider(p string) *Consumer {
	c.provider = p
	return c
}

// ProviderConfig Custom provider configuration (e.g. sarama config, aws credentials)
func (c *Consumer) ProviderConfig(cfg interface{}) *Consumer {
	c.providerConfig = cfg
	return c
}

// Config returns a custom provider configuration (e.g. sarama config, aws credentials)
func (c *Consumer) Config() interface{} {
	return c.providerConfig
}

// Address ip address(es) with its respective port(s) of the Message Broker/Message Queue system cluster
func (c *Consumer) Address(addrs ...string) *Consumer {
	c.cluster = addrs
	return c
}

// Group set of consumers this specific consumer must be with to consume messages in parallel
//
// 	Only available in: Apache Kafka
func (c *Consumer) Group(g string) *Consumer {
	c.group = g
	return c
}

// Publisher pushes the given Message into the Event-Driven ecosystem.
func (c *Consumer) Publisher(p Publisher) *Consumer {
	c.publisher = p
	return c
}

// Handle specific struct Quark will use to send messages
func (c *Consumer) Handle(handler Handler) *Consumer {
	c.handler = handler
	return c
}

// HandleFunc specific func Quark will use to send messages
func (c *Consumer) HandleFunc(handlerFunc HandlerFunc) *Consumer {
	c.handlerFunc = handlerFunc
	return c
}

// TopicString returns every topic registered into the current consumer as string
func (c Consumer) TopicString() string {
	topics := ""
	for i, t := range c.topics {
		if i > 0 {
			topics += ","
		}
		topics += t
	}

	return topics
}
