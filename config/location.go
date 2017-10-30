package config

import "time"

func BogotaLocation() *time.Location {
	bogotaLocation, _ := time.LoadLocation("America/Bogota")
	return bogotaLocation
}
