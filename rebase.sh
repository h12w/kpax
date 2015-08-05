#!/bin/sh

set -e

rm -r -f proto log connector producer
mkdir log
mkdir connector
mkdir producer

wget https://github.com/stealthly/siesta/archive/master.zip
unzip master.zip
mv siesta-master proto
rm master.zip

cd proto
mv connector.go      \
   connector_test.go \
   conn.go           \
   conn_test.go      \
	 ../connector

mv logger.go \
   ../log

mv kafka_producer.go      \
	 kafka_producer_test.go \
	 partitioner.go         \
   record_accumulator.go  \
   network_client.go \
	 selector.go       \
	 selector_test.go  \
   metadata.go       \
	 ../producer

mv utils_test.go util_test.go
cp util.go       \
   util_test.go  \
	 LICENSE       \
	 ../connector

cp LICENSE \
   ../log

#rm README.md
#rm .travis.yml

cd ..

for f in proto/*.go; do
	sed -i.bak 's/package siesta/package proto/' $f
done

for f in log/*.go; do
	sed -i.bak 's/package siesta/package log/' $f
done

for f in connector/*.go; do
	sed -i.bak 's/package siesta/package connector/' $f
done

for f in producer/*.go; do
	sed -i.bak 's/package siesta/package producer/' $f
done

gofmt -w -r 'MetadataResponse -> proto.MetadataResponse' connector/*.go
gofmt -w -r 'FetchResponse -> proto.FetchResponse' connector/*.go
gofmt -w -r 'Broker -> proto.Broker' connector/*.go
gofmt -w -r 'DecodingError -> proto.DecodingError' connector/*.go
gofmt -w -r 'Response -> proto.Response' connector/*.go
gofmt -w -r 'Request -> proto.Request' connector/*.go producer/*.go
gofmt -w -r 'OffsetRequest -> proto.OffsetRequest' connector/*.go
gofmt -w -r 'FetchRequest -> proto.FetchRequest' connector/*.go
gofmt -w -r 'OffsetResponse -> proto.OffsetResponse' connector/*.go
gofmt -w -r 'NewBinaryEncoder-> proto.NewBinaryEncoder' connector/*.go producer/*.go
gofmt -w -r 'NewBinaryDecoder -> proto.NewBinaryDecoder' connector/*.go producer/*.go
gofmt -w -r 'NewOffsetFetchRequest -> proto.NewOffsetFetchRequest' connector/*.go
gofmt -w -r 'OffsetFetchResponse -> proto.OffsetFetchResponse' connector/*.go
gofmt -w -r 'NewMetadataRequest -> proto.NewMetadataRequest' connector/*.go
gofmt -w -r 'ErrReplicaNotAvailable-> proto.ErrReplicaNotAvailable' connector/*.go
gofmt -w -r 'NewConsumerMetadataRequest-> proto.NewConsumerMetadataRequest' connector/*.go
gofmt -w -r 'ConsumerMetadataResponse-> proto.ConsumerMetadataResponse' connector/*.go
gofmt -w -r 'NewOffsetCommitRequest-> proto.NewOffsetCommitRequest' connector/*.go
gofmt -w -r 'OffsetCommitResponse-> proto.OffsetCommitResponse' connector/*.go
gofmt -w -r 'NewRequestHeader-> proto.NewRequestHeader' connector/*.go producer/*.go
gofmt -w -r 'TopicMetadata-> proto.TopicMetadata' connector/*.go
gofmt -w -r 'Message-> proto.Message' producer/*.go connector/*.go
gofmt -w -r 'ErrNoError -> proto.ErrNoError' connector/*.go producer/*.go
gofmt -w -r 'ProduceRequest-> proto.ProduceRequest' producer/*.go
gofmt -w -r 'ProduceResponse-> proto.ProduceResponse' producer/*.go

gofmt -w -r 'Errorf -> log.Errorf' connector/*.go
gofmt -w -r 'Debugf -> log.Debugf' connector/*.go
gofmt -w -r 'Warnf -> log.Warnf' connector/*.go
gofmt -w -r 'Tracef -> log.Tracef' connector/*.go
gofmt -w -r 'Infof -> log.Infof' connector/*.go

gofmt -w -r 'Connector-> connector.Connector' producer/*.go
gofmt -w -r 'BrokerLink-> connector.BrokerLink' producer/*.go

gofmt -w -r 'decodingErr.err -> decodingErr.Error()' producer/*.go
gofmt -w -r 'rawResponseAndError -> RawResponseAndError' connector/*.go
gofmt -w -r 'rawResponseAndError -> connector.RawResponseAndError' producer/*.go
sed -i.bak 's/	bytes \[\]byte/	Bytes \[\]byte/' connector/connector.go
sed -i.bak 's/	link  BrokerLink/	Link  BrokerLink/' connector/connector.go
sed -i.bak 's/	err   error/	Err   error/' connector/connector.go
gofmt -w -r 'response.err -> response.Err' connector/*.go producer/*.go
gofmt -w -r 'response.link -> response.Link' connector/*.go producer/*.go
gofmt -w -r 'response.bytes -> response.Bytes' connector/*.go producer/*.go

# renaming
gofmt -w -r 'NewKafkaProducer -> New' producer/*.go
gofmt -w -r 'RecordMetadata -> RecordError' producer/*.go
gofmt -w -r 'NewDefaultConnector -> New' connector/*.go
gofmt -w -r 'NewConnectorConfig -> NewConfig' connector/*.go

goimports -w log/*.go
goimports -w proto/*.go
goimports -w connector/*.go
goimports -w producer/*.go

cd producer
mv kafka_producer.go      producer.go
mv kafka_producer_test.go producer_test.go
mv record_accumulator.go  accumulator.go
cd ..

rm proto/*.bak
rm log/*.bak
rm connector/*.bak
rm producer/*.bak
