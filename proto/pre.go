package proto

type T interface {
	Marshal() []byte
	Unmarshal([]byte) error
}

func strlen(s string) int16 {
	l := int16(len(s))
	if l == 0 {
		return -1
	}
	return l
}
