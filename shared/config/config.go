package config

import "time"

const (
	Round            int32 = 4
	DefaultGraceDays uint  = 10
)

func DefaultLocation() *time.Location {
	bogotaLocation, _ := time.LoadLocation("America/Bogota")
	return bogotaLocation
}
