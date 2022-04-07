package common

import "time"

// check if start date is before end date
func CheckTimestampIsValid(startTime, endTime time.Time) bool {
	if startTime.Before(endTime) {
		return true
	} else {
		return false
	}
}

func CheckTimestampIsUpToDate(startTime time.Time) bool {
	now := time.Now()
	if now.Before(startTime) {
		return true
	} else {
		return false
	}
}
