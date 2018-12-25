package chatbase

import (
	"time"
)

// TimeStamp returns the current time in UNIX milliseconds
var TimeStamp = func() int64 {
	return time.Now().UnixNano()/1e6
}
