package proto

import (
	"testing"
)

func TestErrorEqual(t *testing.T) {
	var err error
	err = ErrOffsetOutOfRange
	if err != ErrOffsetOutOfRange {
		t.Fatal("error should be equal")
	}
}
