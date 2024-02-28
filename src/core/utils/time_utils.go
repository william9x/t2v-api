package utils

import (
	"time"
)

func MillisToTime(ms int64) (time.Time, error) {
	return time.UnixMilli(ms), nil
}
