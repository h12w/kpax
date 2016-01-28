#!/bin/sh

go build
./gen go  ../BNF.txt   > ../../proto/type_auto.go;    go fmt ../../proto/type_auto.go
./gen gof ../BNF.txt   > ../../proto/marshal_auto.go; go fmt ../../proto/marshal_auto.go
./gen goe ../spec.html > ../../proto/error_auto.go;   go fmt ../../proto/error_auto.go
