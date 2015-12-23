package broker

import (
	"bytes"
	"encoding/binary"

	"github.com/golang/snappy"
)

func decodeSnappy(src []byte) ([]byte, error) {
	if bytes.Equal(src[:8], []byte{130, 83, 78, 65, 80, 80, 89, 0}) {
		result := make([]byte, 0, len(src))
		current := 16
		for current < len(src) {
			size := int(binary.BigEndian.Uint32(src[current : current+4]))
			current += 4
			chunk, err := snappy.Decode(nil, src[current:current+size])
			if err != nil {
				return nil, err
			}
			current += size
			result = append(result, chunk...)
		}
		return result, nil
	}
	return snappy.Decode(nil, src)
}
