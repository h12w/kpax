broker
------

Design principle
----------------
* Modular
* Stateless
* fail fast
  - timeout (SetDealine)
  - release resources
* fault tolerance:
  - try again for broken connection (maxBadConnRetries=2)
  - recover without restart
      - broken connections
      - leader down
  - retry
  - partition expand
  - graceful shutdown
