package quark

import (
	"context"
)

type awsWorker struct {
	id     int
	parent *node
	cfg    AWSConfiguration
}

func (a *awsWorker) SetID(i int) {
	a.id = i
}

func (a *awsWorker) Parent() *node {
	return a.parent
}

func (a *awsWorker) StartJob(ctx context.Context) error {
	return nil
}

func (a *awsWorker) Close() error {
	return nil
}
