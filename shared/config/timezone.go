package config

import "time"

func DefaultLocation() *time.Location {
	bogotaLocation, _ := time.LoadLocation("America/Bogota")
	return bogotaLocation
}
