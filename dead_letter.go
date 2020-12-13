package quark

import (
	"fmt"
)

type DeadLetter struct {
	Topic        string
	Queue        string
	Workers      uint
	Handler      interface{}
	ErrorHandler ErrorHandler
}

func NewDeadLetter(topic, queue string, workers uint, handler interface{}, errorHandler ErrorHandler) *DeadLetter {
	return &DeadLetter{
		Topic:        fmt.Sprintf("dead_letter.%s", topic),
		Queue:        fmt.Sprintf("dead_letter.%s", queue),
		Workers:      workers,
		Handler:      handler,
		ErrorHandler: errorHandler,
	}
}
