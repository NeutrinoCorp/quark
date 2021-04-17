package quark

import "github.com/google/uuid"

// IDFactory Message unique identifier generator
//
//	Default is `google/uuid`
type IDFactory func() string

var defaultIDFactory IDFactory = func() string {
	return uuid.New().String()
}
