package models

import "time"

type Department struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	Name      string     `gorm:"not null"`
	CountryID uint       `gorm:"not null" json:"countryID"`
	Cities    []City     `json:"cities"`
}
