package pkg

import (
	"context"
	"log"
)

type kafkaWorker struct {
	parent *node
	cfg    KafkaConfiguration
}

func (k *kafkaWorker) Parent() *node {
	return k.parent
}

func (k *kafkaWorker) StartJob(ctx context.Context) error {
	log.Print(k.parent.Consumer.topics, k.parent.setDefaultProvider())
	// sarama.NewConsumerGroup(k.parent.setDefaultCluster(), "", k.cfg.Config)
	return nil
}

func (k *kafkaWorker) Close() error {
	log.Print(k.parent.Consumer.topics, "closing worker")
	return nil
}
