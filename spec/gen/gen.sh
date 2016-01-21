#!/bin/sh

go build
./gen go  ../BNF.txt   > ../../broker/type_auto.go;    go fmt ../../broker/type_auto.go
./gen gof ../BNF.txt   > ../../broker/marshal_auto.go; go fmt ../../broker/marshal_auto.go
./gen goe ../spec.html > ../../broker/error_auto.go;   go fmt ../../broker/error_auto.go
