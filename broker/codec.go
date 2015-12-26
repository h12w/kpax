package broker

import (
	"bytes"
	"encoding/binary"

	"github.com/golang/snappy"
	"h12.me/wipro"
)

func (m *Message) Compressed() bool {
	return m.Attributes&0x03 > 0
}

func (m *Message) decompressBytes() (bs []byte, err error) {
	const (
		gzipCodec   = 1
		snappyCodec = 2
	)
	switch m.Attributes & 0x03 {
	case gzipCodec:
		panic("gzip not supported yet")
	case snappyCodec:
		bs, err = decodeSnappy(m.Value)
		if err != nil {
			return nil, err
		}
	}
	return bs, nil
}

func (m *Message) Decompress() (res MessageSet, _ error) {
	bs, err := m.decompressBytes()
	if err != nil {
		return nil, err
	}
	rd := &wipro.Reader{B: bs}
	for rd.Offset < len(rd.B) {
		var m OffsetMessage
		m.Unmarshal(rd)
		if rd.Err != nil {
			return nil, rd.Err
		}
		res = append(res, m)
	}
	return res, nil
}

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
