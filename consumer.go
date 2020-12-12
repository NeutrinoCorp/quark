package quark

// Consumer listen to an specific topic with its own reliability mechanisms (default, retry and dead-letter queues)
type Consumer struct {
	Topic   string
	Queue   string
	Workers int
	Handler interface{}

	Retry      *Retry
	DeadLetter *DeadLetter
}

// NewConsumer creates a new Consumer
func NewConsumer(topic, queue string, workers int, handler interface{},
	retry *Retry, letter *DeadLetter) *Consumer {
	return &Consumer{
		Topic:      topic,
		Queue:      queue,
		Workers:    workers,
		Handler:    handler,
		Retry:      retry,
		DeadLetter: letter,
	}
}
