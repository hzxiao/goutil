package goutil

import "time"

func NowMilli() int64 {
	return time.Now().Local().UnixNano() / int64(time.Millisecond)
}

func NowSec() int64 {
	return time.Now().Local().UnixNano() / int64(time.Second)
}
