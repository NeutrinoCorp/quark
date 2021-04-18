package quark

import (
	"context"
	"time"
)

// Option is a unit of configuration of a Broker
type Option interface {
	apply(*options)
}

type options struct {
	providerConfig         interface{}
	cluster                []string
	errHandler             ErrorHandler
	publisher              Publisher
	eventMux               EventMux
	eventWriter            EventWriter
	poolSize               int
	maxRetries             int
	maxConnRetries         int
	retryBackoff           time.Duration
	connRetryBackoff       time.Duration
	messageIDFactory       IDFactory
	workerFactory          WorkerFactory
	baseMessageSource      string
	baseMessageContentType string
	baseContext            context.Context
}

type clusterOption []string

func (o clusterOption) apply(opts *options) {
	opts.cluster = o
}

// WithCluster defines a set of addresses a Broker will use
func WithCluster(addr ...string) Option {
	return clusterOption(addr)
}

type providerConfigOption struct {
	ProviderConfig interface{}
}

func (o providerConfigOption) apply(opts *options) {
	opts.providerConfig = o.ProviderConfig
}

// WithProviderConfiguration defines an specific Provider configuration
func WithProviderConfiguration(cfg interface{}) Option {
	return providerConfigOption{ProviderConfig: cfg}
}

type errHandlerOption struct {
	Handler ErrorHandler
}

func (o errHandlerOption) apply(opts *options) {
	opts.errHandler = o.Handler
}

// WithErrorHandler defines an error hook/middleware that will be executed when an error occurs inside
// Quark low-level internals
func WithErrorHandler(handler ErrorHandler) Option {
	return errHandlerOption{Handler: handler}
}

type publisherOption struct {
	Publisher Publisher
}

func (o publisherOption) apply(opts *options) {
	opts.publisher = o.Publisher
}

// WithPublisher defines a global Publisher that will be used to push messages to the ecosystem through
// the EventWriter
func WithPublisher(p Publisher) Option {
	return publisherOption{Publisher: p}
}

type eventMuxOption struct {
	Mux EventMux
}

func (o eventMuxOption) apply(opts *options) {
	opts.eventMux = o.Mux
}

// WithEventMux defines the EventMux that will be used for Broker's operations
func WithEventMux(mux EventMux) Option {
	return eventMuxOption{Mux: mux}
}

type eventWriterOption struct {
	EventWriter EventWriter
}

func (o eventWriterOption) apply(opts *options) {
	opts.eventWriter = o.EventWriter
}

// WithEventWriter defines the global event writer
func WithEventWriter(w EventWriter) Option {
	return eventWriterOption{EventWriter: w}
}

type poolSizeOption int

func (o poolSizeOption) apply(opts *options) {
	opts.poolSize = int(o)
}

// WithPoolSize defines the global pool size of total Workers
func WithPoolSize(size int) Option {
	if size <= 0 {
		return poolSizeOption(defaultPoolSize)
	}

	return poolSizeOption(size)
}

type maxRetriesOption int

func (o maxRetriesOption) apply(opts *options) {
	opts.maxRetries = int(o)
}

// WithMaxRetries defines the global maximum ammount of times a publish operations will retry an Event operation
func WithMaxRetries(t int) Option {
	if t <= 0 {
		return maxRetriesOption(defaultMaxRetries)
	}

	return maxRetriesOption(t)
}

type maxConnRetriesOption int

func (o maxConnRetriesOption) apply(opts *options) {
	opts.maxConnRetries = int(o)
}

// WithMaxConnRetries defines the global maximum ammount of times a consumer worker will retry its operations
func WithMaxConnRetries(t int) Option {
	if t <= 0 {
		return maxConnRetriesOption(defaultConnRetries)
	}

	return maxConnRetriesOption(t)
}

type retryBackoffOption time.Duration

func (o retryBackoffOption) apply(opts *options) {
	opts.retryBackoff = time.Duration(o)
}

// WithRetryBackoff defines a time duration an EventWriter will wait to execute a write operation
func WithRetryBackoff(backoff time.Duration) Option {
	if backoff <= 0 {
		return retryBackoffOption(defaultRetryBackoff)
	}
	return retryBackoffOption(backoff)
}

type connRetryBackoffOption time.Duration

func (o connRetryBackoffOption) apply(opts *options) {
	opts.connRetryBackoff = time.Duration(o)
}

// WithConnRetryBackoff defines a time duration a Worker will wait to connect to the specified infrastructure
func WithConnRetryBackoff(backoff time.Duration) Option {
	if backoff <= 0 {
		return connRetryBackoffOption(defaultRetryBackoff)
	}
	return connRetryBackoffOption(backoff)
}

type idFactoryOption struct {
	Factory IDFactory
}

func (o idFactoryOption) apply(opts *options) {
	opts.messageIDFactory = o.Factory
}

// WithMessageIDFactory defines the global Message ID Factory
func WithMessageIDFactory(factory IDFactory) Option {
	return idFactoryOption{Factory: factory}
}

type workerFactoryOption struct {
	Factory WorkerFactory
}

func (o workerFactoryOption) apply(opts *options) {
	opts.workerFactory = o.Factory
}

// WithWorkerFactory defines the global Worker Factory Quark's Supervisor(s) will use to schedule background I/O tasks
func WithWorkerFactory(factory WorkerFactory) Option {
	return workerFactoryOption{Factory: factory}
}

type messageSourceOption string

func (o messageSourceOption) apply(opts *options) {
	opts.baseMessageSource = string(o)
}

// WithBaseMessageSource defines the global base Message Source to comply with the CNCF CloudEvents specification
func WithBaseMessageSource(s string) Option {
	return messageSourceOption(s)
}

type messageTypeOption string

func (o messageTypeOption) apply(opts *options) {
	opts.baseMessageContentType = string(o)
}

// WithBaseMessageContentType defines the global base Message Content Type to comply with the CNCF CloudEvents specification
func WithBaseMessageContentType(s string) Option {
	return messageTypeOption(s)
}

type baseContextOption struct {
	Ctx context.Context
}

func (o baseContextOption) apply(opts *options) {
	opts.baseContext = o.Ctx
}

// WithBaseContext defines the context a Broker will use to executes its operations
func WithBaseContext(ctx context.Context) Option {
	return baseContextOption{Ctx: ctx}
}
