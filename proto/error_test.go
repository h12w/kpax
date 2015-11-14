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
}
