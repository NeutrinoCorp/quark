package pkg

import "context"

type defaultKafkaPublisher struct{}

func (d *defaultKafkaPublisher) Publish(ctx context.Context, message ...*Message) error {
	panic("implement me")
}
