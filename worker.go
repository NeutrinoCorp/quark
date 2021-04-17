package quark

import (
	"context"
)

// Worker is a unit of work a Consumer Supervisor.
//
// It handles all the work and should be running inside a goroutine to enable parallelism.
type Worker interface {
	// SetID inject the current Worker id
	SetID(i int)
	// Parent returns the parent Supervisor
	Parent() *Supervisor
	// StartJob starts an specific work
	//
	// Panics goroutine if any errors is thrown
	StartJob(context.Context) error
	// Close stop all Blocking I/O operations
	Close() error
}

// WorkerFactory is a crucial Broker and/or Consumer component which generates the concrete workers
// Quark will use to consume data
type WorkerFactory func(parent *Supervisor) Worker
