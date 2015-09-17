broker
------

Design principle
----------------
* Modular
* Stateless
* fail fast
  - timeout (SetDealine)
* fault tolerance:
  - automatic handling of broken connections
  - leader down
  - partition expand
  - graceful shutdown
