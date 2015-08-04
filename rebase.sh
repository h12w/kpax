#!/bin/sh

set -e

rm -r -f proto client log
mkdir client
mkdir log

wget https://github.com/stealthly/siesta/archive/master.zip
unzip master.zip
mv siesta-master proto
rm master.zip

cd proto
mv connector.go      \
   connector_test.go \
   conn.go           \
   conn_test.go      \
   metadata.go       \
   network_client.go \
   record_accumulator.go \
	 kafka_producer.go \
	 kafka_producer_test.go \
	 partitioner.go    \
	 selector.go       \
	 selector_test.go  \
	 ../client

mv logger.go \
   ../log

cp util.go       \
   utils_test.go \
	 LICENSE       \
	 ../client

cp LICENSE \
   ../log

rm README.md

cd ..

for f in proto/*.go; do
	sed -i.bak 's/package siesta/package proto/' $f
done


for f in client/*.go; do
	sed -i.bak 's/package siesta/package client/' $f
done

for f in log/*.go; do
	sed -i.bak 's/package siesta/package log/' $f
done

rm proto/*.bak
rm client/*.bak
rm log/*.bak

gofmt -w -r 'MetadataResponse -> proto.MetadataResponse' client/*.go
gofmt -w -r 'FetchResponse -> proto.FetchResponse' client/*.go
gofmt -w -r 'Broker -> proto.Broker' client/*.go
gofmt -w -r 'DecodingError -> proto.DecodingError' client/*.go
gofmt -w -r 'Response -> proto.Response' client/*.go
gofmt -w -r 'Request -> proto.Request' client/*.go
gofmt -w -r 'OffsetRequest -> proto.OffsetRequest' client/*.go
gofmt -w -r 'FetchRequest -> proto.FetchRequest' client/*.go
gofmt -w -r 'OffsetResponse -> proto.OffsetResponse' client/*.go
gofmt -w -r 'NewBinaryDecoder -> proto.NewBinaryDecoder' client/*.go
gofmt -w -r 'NewOffsetFetchRequest -> proto.NewOffsetFetchRequest' client/*.go
gofmt -w -r 'OffsetFetchResponse -> proto.OffsetFetchResponse' client/*.go
gofmt -w -r 'NewMetadataRequest -> proto.NewMetadataRequest' client/*.go
gofmt -w -r 'ErrNoError -> proto.ErrNoError' client/*.go
gofmt -w -r 'ErrReplicaNotAvailable-> proto.ErrReplicaNotAvailable' client/*.go
gofmt -w -r 'NewConsumerMetadataRequest-> proto.NewConsumerMetadataRequest' client/*.go
gofmt -w -r 'ConsumerMetadataResponse-> proto.ConsumerMetadataResponse' client/*.go
gofmt -w -r 'NewOffsetCommitRequest-> proto.NewOffsetCommitRequest' client/*.go
gofmt -w -r 'OffsetCommitResponse-> proto.OffsetCommitResponse' client/*.go
gofmt -w -r 'NewRequestHeader-> proto.NewRequestHeader' client/*.go
gofmt -w -r 'NewBinaryEncoder-> proto.NewBinaryEncoder' client/*.go
gofmt -w -r 'TopicMetadata-> proto.TopicMetadata' client/*.go
gofmt -w -r 'Message-> proto.Message' client/*.go
gofmt -w -r 'ProduceRequest-> proto.ProduceRequest' client/*.go
gofmt -w -r 'ProduceResponse-> proto.ProduceResponse' client/*.go

gofmt -w -r 'Errorf -> log.Errorf' client/*.go
gofmt -w -r 'Debugf -> log.Debugf' client/*.go
gofmt -w -r 'Warnf -> log.Warnf' client/*.go
gofmt -w -r 'Tracef -> log.Tracef' client/*.go
gofmt -w -r 'Infof -> log.Infof' client/*.go

gofmt -w -r 'decodingErr.err -> decodingErr.Error()' client/*.go

goimports -w log/*.go
goimports -w proto/*.go
goimports -w client/*.go

cd client
mv kafka_producer.go      producer.go
mv kafka_producer_test.go producer_test.go
cd ..
