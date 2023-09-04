package scanips

import (
	"fmt"
	"time"
)

func StartTime() string {
	start := time.Now()
	startHour := start.Hour()
	startMin := start.Minute()
	startSec := start.Second()
	startTime := fmt.Sprintf("%d:%d:%d", startHour, startMin, startSec)
	return startTime
}

func EndTime() string {
	end := time.Now()
	endHour := end.Hour()
	endMin := end.Minute()
	endSec := end.Second()
	endTime := fmt.Sprintf("%d:%d:%d", endHour, endMin, endSec)
	return endTime
}
