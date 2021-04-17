package quark

import "errors"

var (
	// ErrBrokerClosed the broker was already closed
	ErrBrokerClosed = errors.New("broker closed")
	// ErrProviderNotValid the given provider is not accepted by Quark
	ErrProviderNotValid = errors.New("provider is not valid")
	// ErrPublisherNotImplemented the given publisher does not have its concrete implementation
	ErrPublisherNotImplemented = errors.New("publisher is not implemented")
	// ErrNotEnoughTopics no topics where found
	ErrNotEnoughTopics = errors.New("not enough topics")
	// ErrNotEnoughHandlers no consumer handler was found
	ErrNotEnoughHandlers = errors.New("not enough handlers")
	// ErrEmptyMessage no message was found
	ErrEmptyMessage = errors.New("message is empty")
	// ErrEmptyCluster the current cluster does not contains any hosts/addresses
	ErrEmptyCluster = errors.New("consumer cluster is empty")
	// ErrRequiredGroup a consumer group is required
	ErrRequiredGroup = errors.New("consumer group is required")
)
