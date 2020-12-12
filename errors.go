package quark

import "errors"

var (
	ErrConsumerAlreadyRegistered = errors.New("subscribers was already registered")
	ErrConsumerHandlerNotValid   = errors.New("consumer handler is not valid")
)
