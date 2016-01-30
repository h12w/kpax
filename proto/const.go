package proto

import (
	"time"
)

const (
	AckNone  = 0
	AckLocal = 1
	AckAll   = -1
)

var (
	Earliest = time.Time{}
	Latest   = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
)
