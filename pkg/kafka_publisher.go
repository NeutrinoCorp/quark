package pkg

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/go-multierror"
)

type defaultKafkaPublisher struct {
	cfg     KafkaConfiguration
	cluster []string
}

func (d *defaultKafkaPublisher) Publish(ctx context.Context, messages ...*Message) error {
	p, err := sarama.NewSyncProducer(d.cluster, d.cfg.Config)
	if err != nil {
		return err
	}
	defer func() {
		err = p.Close()
	}()

	errs := new(multierror.Error)
	for _, msg := range messages {
		errs = multierror.Append(errs, d.sendMessage(ctx, p, msg))
	}
	return errs.ErrorOrNil()
}

func (d *defaultKafkaPublisher) sendMessage(ctx context.Context, p sarama.SyncProducer, msg *Message) error {
	m := MarshalKafkaMessage(msg)
	partition, offset, err := p.SendMessage(m)
	if err != nil {
		return err
	}
	if d.cfg.Producer.OnSent != nil {
		go d.cfg.Producer.OnSent(ctx, m, partition, offset)
	}

	return nil
}
