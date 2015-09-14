h12.me/kafka
============

`h12.me/kafka` is a modular and idiomatic Kafka client in Go:

* **kafka/proto**: Kafka wire protocol (automatically generated types and
  marshal/unmarshal functions from
  [the HTML spec](https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol)).
* **kafka/broker**: client that talks to a single Kafka Broker (concurrent,
  synchronous API with asynchronous request/response IO).
* **kafka/client**: client that talks to a Kafka cluster (metadata management
  & error handling).
* **kafka/producer**: high-level producer (partition, batch and compression strategy).
* **kafka/consumer**: high-level consumer.

Install
-------

```
go get -u h12.me/kafka
```
