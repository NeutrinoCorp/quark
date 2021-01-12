package quark

import "github.com/google/uuid"

// IdGenerator Message unique identifier generator
//
//	Default is `google/uuid`
type IdGenerator func() string

var defaultIdGenerator IdGenerator = func() string {
	return uuid.New().String()
}
