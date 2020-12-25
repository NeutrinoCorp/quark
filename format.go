package quark

import (
	"strconv"
)

const (
	Command     = "command"
	DomainEvent = "domain_event"
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
