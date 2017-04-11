package utils

import (
	"time"
)

var Now = func() int64 {
	return time.Now().UnixNano()
}
