package proto

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	fmt.Println(ErrUnknown)
	fmt.Println(ErrOffsetMetadataTooLargeCode)
	fmt.Println(ErrOffsetsLoadInProgressCode)
	fmt.Println(ErrNotCoordinatorForConsumerCode)
	var err error
	err = ErrOffsetOutOfRange
	if err != ErrOffsetOutOfRange {
		t.Fatal("error should be equal")
	}
}
