h12.me/kafka: a modular and customizable Kafka client
=====================================================

`h12.me/kafka` is a modular and customizable Kafka client package in Go:

* h12.me/kafka/proto: Kafka wire protocol (mostly automatically generated from
  [the HTML spec](https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol)).
* h12.me/kafka/client: a Kafka client with conneciton pooling, bootstrapping and
  error handling
* h12.me/kafka/consumer: high-level consumer
* h12.me/kafka/producer: high-level producer

Install
-------

```
go get -u h12.me/kafka
```
