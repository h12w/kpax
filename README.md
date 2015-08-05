kafka
=====

Forked from github.com/stealthly/siesta.

* kafka/proto: wire protocol only
* kafka/client: client to a Kafka cluster

Docker Cluster for Test
-----------------------

```bash
docker pull wurstmeister/kafka
docker pull wurstmeister/zookeeper

docker-compose up -d
docker-compose scale kafka=3
```
