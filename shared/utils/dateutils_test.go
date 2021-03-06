package utils

import (
	"loans/shared/config"
	"testing"
	"time"
)

func TestDaysSince(t *testing.T) {
	before := time.Date(2017, 8, 1, 0, 0, 0, 0, config.DefaultLocation())
	after := time.Date(2017, 10, 30, 15, 59, 0, 0, config.DefaultLocation())
	days := DaysBetween(after, before)
	daysExpected := 90
	if days != daysExpected {
		t.Fatalf("Expected %v but got %v", daysExpected, days)
	}
}
