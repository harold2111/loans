package models

import (
	"github.com/jinzhu/gorm"
)

type Department struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Country   Country
	CountryID int `gorm:"not null"`
}
