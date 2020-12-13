package brokers

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/go-multierror"
	"github.com/neutrinocorp/quark"
)

type Kafka struct {
	Cluster []string
	Config  *sarama.Config

	nodes map[string]kafkaNode
	mu    *sync.Mutex
	wg    *sync.WaitGroup
}

func NewKafka(cluster []string, config *sarama.Config) *Kafka {
	return &Kafka{
		Cluster: cluster,
		Config:  config,
		nodes:   map[string]kafkaNode{},
		mu:      new(sync.Mutex),
		wg:      new(sync.WaitGroup),
	}
}

func (k *Kafka) Register(consumers ...quark.Consumer) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if err := k.ensureConsumerUniqueness(consumers...); err != nil {
		return err
	}

	for _, consumer := range consumers {
		k.nodes[consumer.Topic] = kafkaNode{
			Cluster:  k.Cluster,
			Config:   k.Config,
			Consumer: consumer,
			Wg:       k.wg,
		}
	}
	return nil
}

func (k Kafka) ensureConsumerUniqueness(consumers ...quark.Consumer) error {
	errs := new(multierror.Error)
	for _, consumer := range consumers {
		if k.IsConsumerRegistered(consumer.Topic) {
			errs = multierror.Append(errs, quark.ErrConsumerAlreadyRegistered)
		}
	}
	return errs.ErrorOrNil()
}

func (k Kafka) IsConsumerRegistered(consumer string) bool {
	_, ok := k.nodes[consumer]
	return ok
}

func (k Kafka) ListenAndServe(ctx context.Context) error {
	totalNodes := k.calculateTotalConsumers()
	k.wg.Add(totalNodes)
	nodeErrChan := make(chan error, totalNodes) // defined chan buffer, each node should have 3 subscribers
	defer close(nodeErrChan)

	k.runNodes(ctx, nodeErrChan)
	select {
	case err := <-nodeErrChan:
		return err
	case <-ctx.Done():
		k.wg.Wait()
		return nil
	}
}

func (k Kafka) calculateTotalConsumers() int {
	total := len(k.nodes)
	for _, n := range k.nodes {
		if n.Consumer.Retry != nil {
			total++
		}
		if n.Consumer.DeadLetter != nil {
			total++
		}
	}
	return total
}

func (k Kafka) runNodes(ctx context.Context, errChan chan error) {
	for _, node := range k.nodes {
		go node.run(ctx, errChan)
	}
}
