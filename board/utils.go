package board

import (
	"fmt"
	"time"
)

func Trace(s string) (string, time.Time) {
	return s, time.Now()
}

func Un(s string, startTime time.Time) {
	endTime := time.Now()
	fmt.Println(s, ":", endTime.Sub(startTime))
}
