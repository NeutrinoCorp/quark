package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/go-multierror"
	"github.com/neutrinocorp/quark"
)

// KafkaPublisher Quark default publisher for Kafka
type KafkaPublisher struct {
	cfg     KafkaConfiguration
	cluster []string
}

// NewKafkaPublisher allocates a new KafkaPublisher
func NewKafkaPublisher(cfg KafkaConfiguration, addrs ...string) *KafkaPublisher {
	return &KafkaPublisher{
		cfg:     cfg,
		cluster: addrs,
	}
}

func (d *KafkaPublisher) Publish(ctx context.Context, messages ...*quark.Message) error {
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

func (d *KafkaPublisher) sendMessage(ctx context.Context, p sarama.SyncProducer, msg *quark.Message) error {
	kafkaMsg := MarshalKafkaMessage(msg)
	partition, offset, err := p.SendMessage(kafkaMsg)
	if err != nil {
		return err
	}
	if d.cfg.Producer.OnSent != nil {
		go d.cfg.Producer.OnSent(ctx, kafkaMsg, partition, offset)
	}

	return nil
}
