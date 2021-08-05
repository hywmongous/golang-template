package values

import "time"

type Time int64

func Now() Time {
	return Time(time.Now().Unix())
}

func (time Time) GetInt64() int64 {
	return int64(time)
}
