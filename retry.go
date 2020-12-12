package quark

import (
	"fmt"
)

type Retry struct {
	Topic        string
	Queue        string
	TotalRetries int
	Workers      int
	Handler      interface{}
}

func NewRetry(topic, queue string, numRetries, workers int, handler interface{}) *Retry {
	return &Retry{
		Topic:        fmt.Sprintf("retry.%s", topic),
		Queue:        fmt.Sprintf("retry.%s", queue),
		TotalRetries: numRetries,
		Workers:      workers,
		Handler:      handler,
	}
}
