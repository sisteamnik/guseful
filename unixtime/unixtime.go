package unixtime

import (
	"time"
)

func Parse(num int64) time.Time {
	t := time.Unix(num/1000000000, num%1000000000)
	t.In(time.Now().Location())
	return t
}
