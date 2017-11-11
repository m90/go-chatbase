package chatbase

import (
	"time"
)

// TimeStamp returns the current time in UNIX millisecond
var TimeStamp = func() int {
	return int(time.Now().Unix()) / int(time.Microsecond)
}
