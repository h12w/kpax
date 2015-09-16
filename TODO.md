broker
------

Design principle
----------------
* Modular
* Stateless
* fail fast
  - SetDeadline
* tolerant special cases:
  - leader down
  - partition expand
  - graceful shutdown
