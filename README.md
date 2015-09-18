h12.me/kafka
============

A modular & idiomatic Kafka client in Go
----------------------------------------

Compatible with Kafka 0.8.2 or later.

Install
-------

```
go get -u h12.me/kafka
```

Design
------

`(-)` means to be done.

### Modular

* **proto**: Kafka Wire Protocol (automatically generated types and
  marshal/unmarshal functions from
  [the HTML spec](https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol)).
* **broker**: client that talks to a single Kafka broker (concurrent,
  synchronous API wraps inside asynchronous request/response IO).
* **client**: client that talks to a Kafka cluster (leader/coordinator management).
* **producer**: fault tolerant high-level producer (batching and partitioning strategy).
* **consumer**: fault tolerant high-level consumer (consumer group and offset commit).

### Error Handling

* fail fast
  + timeout
  + release resources carefully
* fault tolerance
  + retry
    - when a leader down, try another partition
    - try connecting one more time for broken connection (-)
  + recover after broker back online
      - broker reconnect
      - failed partition will be retried after a period of time (producer)
      - metadata reload lazily
  + partition expand (-)
  + graceful shutdown (-)

### Efficiency

* efficiency
  + batching (-)
  + compression (-)
