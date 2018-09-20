package models

import "time"

type Department struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	Name      string     `gorm:"not null"`
	CountryID uint       `gorm:"not null" json:"countryID"`
	Cities    []City     `json:"cities"`
}
