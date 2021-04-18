# :zap: Quark [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Report Status][rep-img]][rep] [![Codebeat][beat-img]][beat] [![Go Version][go-img]][go]
A **_Reliable_** and **_fully customizable_** **event router** for Event-Driven systems written in Go.

Based on reliable mechanisms from companies such as [Uber](https://eng.uber.com/reliable-reprocessing/), 
`Quark` offers an Event Router with a fine-tuned set of tools to ease messaging communication complexity.

Thread-safe processing, parallelism, concurrency and graceful shutdowns are the elemental principles of `Quark`.

Furthermore, `Quark` uses the _[Cloud Native Computing Foundation (CNCF) CloudEvents](https://cloudevents.io/)_ specification to compose messages. `Quark` lets developers use their preferred encoding format _(JSON, Apache Avro, etc.)_ and sets message headers as binary data when possible to reduce computational costs.

Aside basic functionalities, it is worth to mention `Quark` is **_fully customizable_**, so any developer may get the maximum potential out of `Quark`. 

A simple set of examples would be:
- Override the default Event Writer to apply custom resilience mechanisms.
- Increasing a Worker pool size for an specific Consumer process.
- Override the default Publisher (e.g. Apache Kafka) for another provider Publisher (e.g. AWS SNS).
- Set a tracing context as custom Message header to trace an specific event process.

To conclude, `Quark` exposes a friendly API based on Go's idiomatic best practices and the `net/http` + popular HTTP mux (`gorilla/mux`, `gin-gonic/gin`, `labstack/echo`) packages to increase overall usability and productivity.

_More information about the internal low-level architecture may be found [here][quark-arch]._

## Supported Infrastructure
- Apache Kafka
- In Memory*
- Redis Pub/Sub*
- Amazon Web Services Simple Queue Service (SQS)*
- Amazon Web Services Simple Notification Service (SNS)*
- Amazon Web Services Kinesis*
- Amazon Web Services EventBridge*
- Google Cloud Platform Pub/Sub*
- Microsoft Azure Service Bus*
- NATS*
- RabbitMQ*

_* to be implemented_

## Installation

Since `Quark` uses Go submodules to decompose specific depenencies for providers, it is required to install concrete implementations _(Apache Kafka, In memory, Redis, ...)_ manually. One may install these using the following command.

`go get github.com/neutrinocorp/quark/bus/YOUR_PROVIDER`

If one wants to develop its own custom implementations, it is required to install `Quark` core library. It can be done running the following command.

`go get github.com/neutrinocorp/quark`

_Note that `Quark` only supports the two most recent minor versions of Go._

## Quick Start

Before we set up our consumers, we must define our `Broker` and its required configuration to work as desired.

The following example demonstrates how to set up an _Apache Kafka_ `Broker` with an error handler (hook).

```go
// Create error hook
customErrHandler := func(ctx context.Context, err error) {
  log.Print(err)
}

// ...

// Create broker
b := kafka.NewKafkaBroker(
	// ... A Shopify/sarama configuration,
	quark.WithCluster("localhost:9092", "localhost:9093"),
	quark.WithBaseMessageSource("https://neutrinocorp.org/cloudevents"),
	quark.WithBaseMessageContentType("application/cloudevents+json"),
	quark.WithErrorHandler(customErrHandler))
```

### Listen to a Topic

Quark is very straight forward as is based on the `net/http` and well known Go HTTP mux packages.

This example demonstrates how to listen to a Topic using the `Broker.Topic()` method.

If no pool-size was specified, `Quark` will set up to 5 `workers` for each Consumer.

```go
b.Topic("chat.1").HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  // ...
  return true
})
```

### Publish to Topic(s)

Internally, Quark uses the `EventWriter` and `Publisher` components to propagate an event into the provider infrastructure.

Moreover, Quark lets developers publish a message to multiple Topics (fan-out).

This can be done by calling the `EventWriter's Write()/WriteMessage()` methods.

```go
b.Topic("chat.1").HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  // ...
  _, _ = w.Write(e.Context, encodedMsgBody, "chat.2", "chat.3") // returns how many messages were published
  // or
  // msg is a quark.Message struct, writer will use Message.Type/Topic attribute to publish
  _, _ = w.WriteMessage(e.Context, msgA, msgB)
  return true
})
```

### Increase/Decrease Worker pool for a Consumer process

Quark parallelize message-processing jobs by creating a pool of `Worker(s)` for each Consumer using goroutines.

The pool size can be defined to an specific Consumer calling the `quark.WithPoolSize()/Consumer.PoolSize()` method/function.

```go
b.Topic("chat.1").PoolSize(10).HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  // ...
  return true
})
```

### Retry an event process

Quark is based on _reliable mechanisms_ such as _retry-exponential+jitter_ backoff and sending _poison messages_ to Dead-Letter Queues (DLQ) strategies.

To customize these mechanisms, the developer may use the `quark.WithMaxRetries()/Consumer.MaxRetries()` and `quark.WithRetryBackoff()/Consumer.RetryBackoff()` methods/functions.

_These strategies are implemented by default on the `defaultEventWriter` component._

```go
b.Topic("cosmos.payments").MaxRetries(3).RetryBackoff(time.Second*3).HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  // ... something failed in our processing
  err := w.WriteRetry(e.Context, e.Body)
  if errors.Is(err, quark.ErrMessageRedeliveredTooMuch) {
       	// calling Write() will set the Message re-delivery delta to 0
	_, _ = w.Write(e.Context, e.Body.Data, "dlq.chat.1")
  }
  return true
})
```

### Failed event processing

If a message processing fails, `Quark` will use _**Acknowledgement**_ mechanisms if available.

This can be done by returning a `false` value from the event handler.

_* Only available for specific providers._

```go
b.Topic("cosmos.user_registered").HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  // ...
  return true // this indicates if the consumer should mark the message or not (Ack or NAck)
})
```

### Start Broker and Graceful Shutdown

To conclude, after setting up all of our consumers, the developer must start the `Broker` component to execute background jobs from registered `Consumer(s)`.

The developer should not forget to shutdown gracefully the `Broker` _-like an `net/http` server-_.

```go
// graceful shutdown
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt)
go func() {
	if err := b.ListenAndServe(); err != nil && err != quark.ErrBrokerClosed {
		log.Fatal(err)
	}
}()

<-stop

log.Printf("stopping %d supervisor(s) and %d worker(s)", b.ActiveSupervisors(), b.ActiveWorkers())
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()

if err := b.Shutdown(ctx); err != nil {
	log.Fatal(err)
}

log.Print(b.ActiveSupervisors(), b.ActiveWorkers()) // should be 0,0
```

## Advanced techniques

### Grouping Consumer jobs

When processing in parallel, every Worker in a Consumer pool will read from a Queue/Offset independently.

Even though this is intended by Quark, the developer might want to balance processing load from the Worker(s) by treating the Consumer pool as a whole.

This can be done calling the `Consumer.Group()` method.

_* Only available for specific providers (e.g. Apache Kafka)._

```go
b.Topic("chat.1").Group("awesome-group").HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  log.Print(e.Topic, e.RawValue)
  // ...
  return true
})
```

### Listening N-Topics within a single Consumer (Fan-in)

A Quark Consumer accepts up to N topics by default.

This feature can be implemented by calling the `Consumer.Topics()` method.

```go
b.Topics("chat.0", "chat.1").HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  log.Print(e.Topic, e.RawValue)
  // ...
  return true
})
```

### Event header read and manipulation

Like HTTP, Quark defines a set of headers for each Event and decodes/encodes them by default.

These headers may contain useful metadata from the current Broker, Consumer and provider (e.g. an Apache Kafka Offset or Partition).

Moreover, Quark lets developers read or manipulate these headers. Thus, modified headers will be published when EventWriter's Write methods are called.

This can be done by calling the `EventWriter.Get()` and `EventWriter.Set()` methods.

```go
b.Topic("chat.1").HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
  // ...
  partition := w.Header().Get(kafka.HeaderKafkaPartition)
  w.Header().Set(quark.HeaderMessageDataContentType, "application/avro")
  _, _ = w.Write(e.Context, e.Body.Data, "dlq.chat.1") // will use new Content-Type header
  return true
})
```

### Using a different Publisher for a Consumer process

As part of the _fully customizable_ principle, a Quark Consumer may use a different Publisher component if desired.

This feature can be implemented by using a different provider `Publisher` implementation and by calling the `Consumer.Publisher()` method.

```go
// on quark/bus/aws package

type SNSPublisher struct{}

func (a SNSPublisher) Publish(ctx context.Context, msgs ...*quark.Message) error {
	// ...
	return nil
}

// on developer application

// ...

// Listening from Google Cloud Platform Pub/Sub

b.Topic("alex.trades").Publisher(aws.SNSPublisher{}).
  HandleFunc(func(w quark.EventWriter, e *quark.Event) bool {
    // Write() will publish the message to Amazon Web Services Simple Notification Service (SNS)
    _, _ = w.Write(e.Context, []byte("alex has traded in a new index fund"),
      "aws.alex.trades", "aws.analytics.trades")
    return true
  })
```

See the [documentation][doc], [examples][examples] and [FAQ](FAQ.md) for more details.

## Performance

As measured by its own [benchmarking suite][benchmarking_suite], not only is `Quark` more performant
than comparable messaging processing packages. Like all benchmarks, take these with a grain of salt.<sup
id="anchor-versions">[1](#footnote-versions)</sup>

## Maintenance
This library is currently maintained by
- [maestre3d][maintainer]

## Development Status: Alpha

All APIs are under development, yet no breaking changes will be made in the 1.x.x series
of releases. Users of semver-aware dependency management systems should pin
`Quark` to `^1`.

## Contributing

We encourage and support an active, healthy community of contributors &mdash;
including you! Details are in the [contribution guide](CONTRIBUTING.md) and
the [code of conduct](CODE_OF_CONDUCT.md). The Quark maintainers keep an eye on
issues and pull requests, but you can also report any negative conduct to
**oss-conduct@neutrinocorp.org**. That email list is a private, safe space; even the Quark
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
[quark-arch]: https://github.com/neutrinocorp/quark/tree/master/docs/quark-arch.png
[ci-img]: https://github.com/neutrinocorp/quark/workflows/Go/badge.svg?branch=master
[ci]: https://github.com/NeutrinoCorp/quark/actions
[go-img]: https://img.shields.io/github/go-mod/go-version/NeutrinoCorp/quark?style=square
[go]: https://github.com/NeutrinoCorp/quark/blob/master/go.mod
[rep-img]: https://goreportcard.com/badge/github.com/neutrinocorp/quark
[rep]: https://goreportcard.com/report/github.com/neutrinocorp/quark
[cov-img]: https://codecov.io/gh/NeutrinoCorp/quark/branch/master/graph/badge.svg
[beat-img]: https://codebeat.co/badges/416103dd-8b2a-463e-83fb-5c438c2565ac
[beat]: https://codebeat.co/projects/github-com-neutrinocorp-quark-master
[cov]: https://codecov.io/gh/NeutrinoCorp/quark
[benchmarking_suite]: https://github.com/neutrinocrp/quark/tree/master/benchmarks
[benchmarks/go.mod]: https://github.com/neutrinocorp/quark/blob/master/benchmarks/go.mod
[maintainer]: https://github.com/maestre3d
