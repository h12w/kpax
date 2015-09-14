h12.me/kafka: a modular & idiomatic Kafka client in Go
======================================================

`h12.me/kafka` is a modular and idiomatic Kafka client in Go:

* h12.me/kafka/proto: Kafka wire protocol (automatically generated types and
  marshal/unmarshal functions from
  [the HTML spec](https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol)).
* h12.me/kafka/broker: client that talks to a single Kafka Broker (concurrent,
  synchronous API with asynchronous request/response IO).
* h12.me/kafka/client: client that talks to a Kafka cluster (metadata management
  & error handling).
* h12.me/kafka/producer: high-level producer (partition, batch and compression strategy).
* h12.me/kafka/consumer: high-level consumer.

Install
-------

```
go get -u h12.me/kafka
```
