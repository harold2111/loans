package config

import "time"

const (
	Round                                  = 4
	DaysBeforeEndDateToConsiderateDue      = -15
	DaysAfterEndDateToConsiderateInArrears = 5
)

func DefaultLocation() *time.Location {
	bogotaLocation, _ := time.LoadLocation("America/Bogota")
	return bogotaLocation
}
