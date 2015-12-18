h12.me/kafka
============

[![GoDoc](https://godoc.org/h12.me/kafka?status.svg)](https://godoc.org/h12.me/kafka)

A modular & idiomatic Kafka client in Go
----------------------------------------

Install
-------

```
go get -u h12.me/kafka
```

Design
------

The client is built on top of Kafka Wire Protocol (i.e. low-level API). The protocol related types & marshal/unmarshal functions are automatically generated by [wipro](https://github.com/h12w/wipro) from [the HTML spec](https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol)).

`(-)` means to be done.

### Sub packages

* **broker**: client that talks to a single Kafka broker (concurrent,
  synchronous API wraps inside asynchronous request/response IO)
* **cluster**: client that talks to a Kafka cluster (leader/coordinator management)
* **producer**: fault tolerant high-level producer (batching and partitioning strategy)
* **consumer**: fault tolerant high-level consumer (consumer group and offset commit)
* **log**: replaceable global logger
* **kafpro**: command line tool to query Kafka wire API

### Compatibility

Compatible with Kafka Server 0.8.2.

### Error Handling

* broker
  + fail fast: timeout returns error immediately
  + release resources carefully
  + reconnect when requested next time
* client
  + metadata reload lazily (only when a leader/coordinator cannot be found in cache)
  + leader/coordinator should be deleted on error
* producer
  + fail over to another partition
  + failed partition will be retried again after a period of time
  + partition expand (-)
* consumer
  + just loop & wait on error
  + partition expand (-)
* graceful shutdown (-)

### Efficiency

* efficiency
  + batching
    - consumer response
    - consumer request (-)
    - producer (-)
  + compression (-)
