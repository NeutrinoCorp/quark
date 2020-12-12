package quark

import (
	"fmt"
)

type DeadLetter struct {
	Topic   string
	Queue   string
	Workers int
	Handler interface{}
}

func NewDeadLetter(topic, queue string, workers int, handler interface{}) *DeadLetter {
	return &DeadLetter{
		Topic:   fmt.Sprintf("dead_letter.%s", topic),
		Queue:   fmt.Sprintf("dead_letter.%s", queue),
		Workers: workers,
		Handler: handler,
	}
}
