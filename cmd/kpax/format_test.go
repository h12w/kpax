package main

import (
	"testing"

	"gopkg.in/vmihailenco/msgpack.v2"
)

func TestUnmarshalMsgPack(t *testing.T) {
	type S struct {
		T struct {
			I int
		}
	}
	var s S
	s.T.I = 1
	buf, err := msgpack.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"T":{"I":1}}`
	actual := MsgPackFormat.Sprint(buf)
	if actual != expected {
		t.Fatalf("expect %s but got %s", expected, actual)
	}
}
