
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
    create_snappy_message([("---start---", None), ("xxxxxxxxxxxxxxxxx", None)]),
    create_snappy_message([("ccccc", None), ("---end---", None)]),
]

produce = ProduceRequestPayload("topic1", 1, messages=messages)

correlationID = 3
req = KafkaProtocol.encode_produce_request("clientID", correlationID, payloads=[produce], acks=1, timeout=1000)

print "[]byte{" +''.join( [ "0x%02x, " % ord( x ) for x in req ] ).strip() + "}"

# python test.py | hexdump -C
import subprocess
task = subprocess.Popen("hexdump -C", shell=True, stdin=subprocess.PIPE)
task.stdin.write(req)
task.stdin.flush()
task.stdin.close()
task.wait()
