h12.me/kafka: a modular Kafka client in Go
==========================================

`h12.me/kafka` is a modular Kafka client in Go:

* h12.me/kafka/proto: Kafka wire protocol (mostly automatically generated from
  [the HTML spec](https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol)).
* h12.me/kafka/broker: client that talks to a single Kafka Broker.
* h12.me/kafka/client: client that talks to a Kafka cluster.
* h12.me/kafka/consumer: high-level consumer.
* h12.me/kafka/producer: high-level producer.

Install
-------

```
go get -u h12.me/kafka
```
