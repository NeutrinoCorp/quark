package quark

import (
	"strconv"
)

const (
	// Command Message type. Used when dispatching an asynchronous command (CQRS); in other words, this type is widely
	// used when asynchronous processing in an external system is required (e.g. processing a video in a serverless
	// function triggered by a service or simple synchronous REST API)
	Command = "command"
	// DomainEvent Message type. Used when something has happened in an specific aggregate and it must propagate
	// its side-effects in the entire event-driven ecosystem
	DomainEvent = "event"
)

// FormatTopicName forms an Async API topic name
//	format e.g. "organization.service.version.kind.entity.action"
func FormatTopicName(organization, service, kind, entity, action string, version int) string {
	return organization + "." + service + "." + strconv.Itoa(version) + "." + kind + "." + entity + "." + action
}

// FormatQueueName forms an Async API queue name
//	format e.g. "service.entity.action_on_event
func FormatQueueName(service, entity, action, event string) string {
	return service + "." + entity + "." + action + "_on_" + event
}
