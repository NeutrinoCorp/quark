# :zap: Quark [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Report Status][rep-img]][rep] [![Go Version][go-img]][go]
Reliable Event-Driven mechanisms for reactive ecosystems written in Go.

Based on reliable mechanisms from companies such as [Uber](https://eng.uber.com/reliable-reprocessing/) who serve 1M+ requests per-hour, 
Quark offers a fine-tuned set of tools to ease overall complexity.

In Addition, `Quark` _fans-out processes per-consumer to_ **parallelize blocking I/O** _tasks_ (as consuming from a queue/topic would be) and isolate them. Thread-safe and graceful shutdown are a very important part of Quark.

Furthermore, `Quark` stores specific data _(e.g. event id, correlation id, span context, etc)_ into messages headers in a binary format to ease disk consumption on the infrastructure, and it lets users use their own encoding preferred library. _(JSON, Apache Avro, etc.)_

Therefore, `Quark` lets developers take advantage of those mechanisms with its default configuration and a `gorilla/mux`-like router/mux to keep them in ease and get benefits without complex configurations and handling. You can either choose use _global configurations_ specified in the broker or use an _specific configuration_ for a specific consumer.

_We currently offer Apache Kafka and AWS SNS/SQS implementations, yet we will be adding more implementations such as GCP PubSub, RabbitMQ and NATS in a near future._

## Installation

`go get github.com/neutrinocorp/quark`

Note that quark only supports the two most recent minor versions of Go.

## Quick Start

Before we set up our consumers, we must define our `Broker` and its required configuration to work as desired.

The following example demonstrates how to set up an _Apache Kafka_ `Broker` with an error handler (hook).

_When using Apache Kafka, `Shopify/sarama` package is required as we rely on its mechanisms._

```go
// Create broker
b := quark.NewKafkaBroker(newSaramaCfg(), "localhost:9092")

b.ErrorHandler = func(ctx context.Context, err error) {
  log.Print(err)
}
```

Quark is very straight forward as is based on the `net/http` and `gorilla/mux` packages.
This example demonstrates how to listen to an asynchronous topic using the `Topic` function.

If no pool-size was specified, `Quark` will set up to 5 `workers` per-consumer node.

```go
b.Topic("chat.1").HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  log.Print(e.Topic, e.RawValue)
  // publish messages to given topics
  _, _ = w.Write(e.Context, e.RawValue, "chat.2", "chat.3") // returns how many messages were published
  return true // this indicates if the consumer should mark the message or not (Ack or NAck)
})
```

Quark parallelize consumers tasks into a pool of `workers` using goroutines and executes a graceful shutdown by default. 

The pool size can be defined by the user with a `PoolSize` attribute.

```go
b.Topic("chat.1").PoolSize(10).HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  log.Print(e.Topic, e.RawValue)
  // publish messages to given topics
  _, _ = w.Write(e.Context, e.RawValue, "chat.2", "chat.3")
  return true
})
```

Quark is based on _reliable mechanisms_. To make use of them, one needs to specify on either the `Broker` or on a specific topic.

This method relies on `defaultEventWriter` as it contains preconfigured reliable mechanisms to avoid message loops and more functionalities.

```go
b.Topic("cosmos.payments").MaxRetries(3).RetryBackoff(time.Second*3).HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  // ... something failed in our processing
  if e.Body.Metadata.RedeliveryCount > 3 {
  	// executed if our message has been processed too much, send to Dead Letter Queue
  	w.Header().Set(quark.HeaderMessageRedeliveryCount, 0)
  	_, _ = w.Write(e.Context, e.RawValue, "dlq.cosmos.payment")
	return true
  }
  
  // publish messages to retry queue, message loop will be avoided by defaultEventWriter
  e.Body.Metadata.RedeliveryCount++
  w.Header().Set(quark.HeaderMessageRedeliveryCount, strconv.Itoa(e.Body.Metadata.RedeliveryCount))
  _, _ = w.Write(e.Context, e.RawValue, "retry.cosmos.payment")
  return true
})
```

To conclude, after setting up all of our consumers, we must start the `Broker` up to trigger and rise all the specified `Consumer`.

Don't forget to graceful shutdown as if you were shutting down a `net/http` server.

```go
// graceful shutdown
stop := make(chan os.Signal)
signal.Notify(stop, os.Interrupt)
go func() {
  if err = b.ListenAndServe(); err != nil && err != quark.ErrBrokerClosed {
    log.Fatal(err)
  }
}()

<-stop

log.Printf("stopping %d nodes and %d workers", b.RunningNodes(), b.RunningWorkers())
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()

if err = b.Shutdown(ctx); err != nil {
  log.Fatal(err)
}

log.Print(b.RunningNodes(), b.RunningWorkers()) // should be 0,0
```

### Advanced techniques

**Grouping consumer nodes**

One might want to take advantage of specific features like _Apache Kafka's_ `Consumer Group` as this might help to process messages one-at-the-time.

This can be done using the default `mux` by setting the `Group` attribute.

By default, Quark uses `Partition Consumer` when using Apache Kafka.

```go
b.Topic("chat.1").Group("awesome-group").HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  log.Print(e.Topic, e.RawValue)
  // publish messages to given topics
  _, _ = w.Write(e.Context, e.RawValue, "chat.2", "chat.3")
  return true
})
```

**Fanning-in messages/queues into a single consumer**

The following example demonstrates how to do the current case using the `mux` and the `Topics` function.

_When fan-in is configured, the `Consumer` must be inside a `Group`_

```go
b.Topics("chat.0", "chat.1").Group("chat-group").HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  log.Print(e.Topic, e.RawValue)
  // publish messages to given topics
  _, _ = w.Write(e.Context, e.RawValue, "chat.2", "chat.3")
  return true
})
```

**Custom publisher per-consumer**

Say you were listening topics from Kafka, yet you want to publish the output into AWS SNS instead Kafka (specified in the global configuration).

The following example demonstrates how to tackle the previous scenario with Quark.

Therefore, the use of `Group` is crucial here since `Partition Consumer` is treated as single unit of processing, and it would publish the message N-times (the pool size since `Consumer` workers are running in parallel).

```go
type AWSPublisher struct{}

func (a AWSPublisher) Publish(ctx context.Context, msgs ...*quark.Message) error {
	for _, msg := range msgs {
		log.Printf("publishing - message: %s", msg.Kind)
	}
	return nil
}

// ...

b.Topic("alex.trades").Group("alex.trades").Publisher(AWSPublisher{}).
  HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
    _, _ = w.Write(e.Context, []byte("alex has traded in a new index fund"),
      "aws.alex.trades", "aws.analytics.trades")
    return true
  })
```

See the [documentation][doc], [examples][examples] and [FAQ](FAQ.md) for more details.

## Performance

As measured by its own [benchmarking suite][], not only is quark more performant
than comparable messaging processors packages. Like all benchmarks, take these with a grain of salt.<sup
id="anchor-versions">[1](#footnote-versions)</sup>

## Supported Infrastructure
- Apache Kafka
- Amazon Web Services SNS and SQS (Separated or using both)

## Maintenance
This library is currently maintained by
- [maestre3d][maintainer]

## Development Status: Alpha

All APIs are under development, no breaking changes will be made in the 1.x series
of releases. Users of semver-aware dependency management systems should pin
quark to `^1`.

## Contributing

We encourage and support an active, healthy community of contributors &mdash;
including you! Details are in the [contribution guide](CONTRIBUTING.md) and
the [code of conduct](CODE_OF_CONDUCT.md). The quark maintainers keep an eye on
issues and pull requests, but you can also report any negative conduct to
oss-conduct@neutrinocorp.org. That email list is a private, safe space; even the zap
maintainers don't have access, so don't hesitate to hold us to a high
standard.

<hr>

Released under the [MIT License](LICENSE).

<sup id="footnote-versions">1</sup> In particular, keep in mind that we may be
benchmarking against slightly older versions of other packages. Versions are
pinned in the [benchmarks/go.mod][] file. [↩](#anchor-versions)

[doc-img]: https://pkg.go.dev/badge/github.com/neutrinocorp/quark
[examples]: https://github.com/neutrinocorp/quark/tree/master/examples
[doc]: https://pkg.go.dev/github.com/neutrinocorp/quark
[docs]: https://github.com/neutrinocorp/quark/tree/master/docs
[ci-img]: https://github.com/neutrinocorp/quark/workflows/Go/badge.svg?branch=master
[ci]: https://github.com/NeutrinoCorp/quark/actions
[go-img]: https://img.shields.io/github/go-mod/go-version/NeutrinoCorp/quark?style=square
[go]: https://github.com/NeutrinoCorp/quark/blob/master/go.mod
[rep-img]: https://goreportcard.com/badge/github.com/neutrinocorp/quark
[rep]: https://goreportcard.com/report/github.com/neutrinocorp/quark
[cov-img]: https://codecov.io/gh/NeutrinoCorp/quark/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/NeutrinoCorp/quark
[benchmarking suite]: https://github.com/neutrinocrp/quark/tree/master/benchmarks
[benchmarks/go.mod]: https://github.com/neutrinocorp/quark/blob/master/benchmarks/go.mod
[maintainer]: https://github.com/maestre3d
