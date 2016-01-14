
# generate test cases with:
#     https://github.com/dpkp/kafka-python
#
# sudo pip install kafka-python

from collections import namedtuple
import functools
from kafka.protocol import KafkaProtocol

from kafka import (
    SimpleProducer, KeyedProducer,
    create_message, create_gzip_message, create_snappy_message,
    RoundRobinPartitioner, HashedPartitioner
)

ProduceRequestPayload = namedtuple("ProduceRequestPayload",
    ["topic", "partition", "messages"])

messages = [
    create_snappy_message([("Snappy 1 %d" % i, None) for i in range(1)]),
    create_snappy_message([("Snappy 2 %d" % i, None) for i in range(2)]),
]

produce = ProduceRequestPayload("topic1", 1, messages=messages)

req = KafkaProtocol.encode_produce_request("abc", 2, payloads=[produce], acks=1, timeout=1000)
print len(req)
print "[]byte{" +''.join( [ "0x%02x, " % ord( x ) for x in req ] ).strip() + "}"
