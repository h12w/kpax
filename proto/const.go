package proto

import (
	"time"
)

type ProduceAckType int16

const (
	AckNone  ProduceAckType = 0
	AckLocal ProduceAckType = 1
	AckAll   ProduceAckType = -1
)

var (
	Earliest = time.Time{}
	Latest   = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
)
