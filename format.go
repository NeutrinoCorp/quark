package quark

import "fmt"

// FormatTopicName forms an Async API topic name
//	format e.g. "organization.service.version.kind.entity.action"
func FormatTopicName(organization, service, kind, entity, action string, version int) string {
	return fmt.Sprintf("%s.%s.%d.%s.%s.%s", organization, service, version, kind, entity, action)
}

// FormatQueueName forms an Async API queue name
//	format e.g. "service.entity.action_on_event
func FormatQueueName(service, entity, action, event string) string {
	return fmt.Sprintf("%s.%s.%s_on_%s", service, entity, action, event)
}
