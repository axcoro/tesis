package board

import (
	"fmt"
	"time"
)

func trace(s string) (string, time.Time) {
	// log.Println("START:", s)
	return s, time.Now()
}

func un(s string, startTime time.Time) {
	endTime := time.Now()
	fmt.Println(s, ":", endTime.Sub(startTime))
}
