package pkg

import "errors"

var (
	ErrBrokerClosed            = errors.New("broker closed")
	ErrProviderNotValid        = errors.New("provider is not valid")
	ErrPublisherNotImplemented = errors.New("publisher is not implemented")
	ErrNotEnoughTopics         = errors.New("not enough topics")
	ErrEmptyCluster            = errors.New("consumer cluster is empty")
)
