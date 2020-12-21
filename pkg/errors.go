package pkg

import "errors"

var (
	ErrProviderNotValid  = errors.New("provider is not valid")
	ErrPublisherNotFound = errors.New("publisher was not found")
)
