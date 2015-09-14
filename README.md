h12.me/kafka: a modular Kafka client in Go
==========================================

`h12.me/kafka` is a modular Kafka client in Go:

* h12.me/kafka/proto: Kafka wire protocol (mostly automatically generated from
  [the HTML spec](https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol)).
* h12.me/kafka/broker: client that talks to a single Kafka Broker (async request/response wrapped in a sync API).
* h12.me/kafka/client: client that talks to a Kafka cluster (metadata management).
* h12.me/kafka/producer: high-level producer (batch strategy, compression, etc).
* h12.me/kafka/consumer: high-level consumer.

Install
-------

```
go get -u h12.me/kafka
```
