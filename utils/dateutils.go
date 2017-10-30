package utils

import "time"

func DaysSince(since time.Time) int {
	d := 24 * time.Hour
	sinceUTC := since.In(time.UTC).Truncate(d)
	return int(time.Since(sinceUTC).Hours() / 24)
}

func AddMothToTimeUtil(startTime time.Time, monthToAdd int) time.Time {
	endTime := startTime.AddDate(0, monthToAdd, 0)
	endTimeWithLastMothDay := time.Date(endTime.Year(), endTime.Month(), 0,
		endTime.Hour(), endTime.Minute(), endTime.Second(), endTime.Nanosecond(), endTime.Location())
	if startTime.Day() > endTimeWithLastMothDay.Day() {
		return endTimeWithLastMothDay
	}
	return endTime
}
