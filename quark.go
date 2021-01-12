// quark Reliable Event-Driven mechanisms for reactive ecosystems written in Go.
//
// Based on reliable mechanisms from companies such as Uber who serve 1M+ requests per-hour, Quark offers a
// fine-tuned set of tools to ease overall complexity.
//
// In Addition, Quark fans-out processes per-consumer to parallelize blocking I/O tasks
// (as consuming from a queue/topic would be) and isolate them. Thread-safe and graceful shutdown are a very important
// part of Quark.
//
// Furthermore, Quark stores specific data (e.g. event id, correlation id, span context, etc) into messages headers in
// binary format to ease disk consumption on the infrastructure and it lets users use their own encoding preferred
// library. (JSON, Apache Avro, etc.)
//
// Therefore, Quark lets developers take advantage of those mechanisms with its default configuration and a
// gorilla/mux-like router/mux to keep them in ease and get benefits without complex configurations and handling.
// You can either choose use global configurations specified in the broker or use an specific configuration for
// an specific consumer.
package quark
