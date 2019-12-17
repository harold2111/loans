package utils

import (
	"time"
)

//fixed := time.FixedZone("-05:00", 0) get time in utc witout fix timezone

func DaysBetween(before, after time.Time) int {
	beforeFixed := FixTimeToZeroHours(before)
	afterFixed := FixTimeToZeroHours(after)
	diff := afterFixed.Sub(beforeFixed)
	hours := diff.Hours()
	days := hours / 24
	return int(days)
}

func FixTimeToZeroHours(timeToFix time.Time) time.Time {
	return time.Date(timeToFix.Year(), timeToFix.Month(), timeToFix.Day(), 0, 0, 0, 0, timeToFix.Location())
}

func AddMothToTimeForPayment(startTime time.Time, monthToAdd int) time.Time {
	endTime := startTime.AddDate(0, monthToAdd, 0)
	endTime = time.Date(endTime.Year(), endTime.Month(), startTime.Day(),
		endTime.Hour(), endTime.Minute(), endTime.Second(), endTime.Nanosecond(), endTime.Location())
	endTimeWithLastMothDay := time.Date(endTime.Year(), endTime.Month(), 0,
		endTime.Hour(), endTime.Minute(), endTime.Second(), endTime.Nanosecond(), endTime.Location())
	if startTime.Day() > endTimeWithLastMothDay.Day() {
		return endTimeWithLastMothDay
	}
	return endTime
}
