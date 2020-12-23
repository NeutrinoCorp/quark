package pkg

import (
	"context"
	"log"
)

type awsWorker struct {
	parent *node
	cfg    AWSConfiguration
}

func (a *awsWorker) Parent() *node {
	return a.parent
}

func (a *awsWorker) StartJob(ctx context.Context) error {
	log.Print(a.parent.Consumer.topics, a.parent.setDefaultProvider())
	return nil
}

func (a *awsWorker) Close() error {
	log.Print(a.parent.Consumer.topics, "closing worker")
	return nil
}
