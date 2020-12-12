package quark

import "context"

type ErrorHandler func(ctx context.Context, topic, queue string, err error)
