package quark

import (
	"fmt"
)

type Retry struct {
	Topic        string
	Queue        string
	TotalRetries uint
	Workers      uint
	Handler      interface{}
	ErrorHandler ErrorHandler
}

func NewRetry(topic, queue string, numRetries, workers uint, handler interface{}, errorHandler ErrorHandler) *Retry {
	return &Retry{
		Topic:        fmt.Sprintf("retry.%s", topic),
		Queue:        fmt.Sprintf("retry.%s", queue),
		TotalRetries: numRetries,
		Workers:      workers,
		Handler:      handler,
		ErrorHandler: errorHandler,
	}
}
