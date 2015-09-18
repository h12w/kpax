TODO
====

Design principle
----------------

`(-)` means to be done.

### Modular

proto: Wire Protocol marshal/unmarshal
broker: single broker connection
client: cluster connection, leader/coordinator management
producer: batching, fault tolerance
consumer: batching, fault tolerance

### Error Handling

* fail fast
  + timeout (SetDealine)
  + release resources carefully
* fault tolerance:
  + retry
    - redirect produce request to another partition (-)
    - try again for broken connection (-)
  + recover after broker back online without restart
      - broker reconnect
      - metadata reload
  + partition expand (-)
  + graceful shutdown (-)

### Efficiency

* efficiency
  + batching (-)
  + compression (-)
