package models

import (
	"github.com/jinzhu/gorm"
)

type Address struct {
	gorm.Model
	ClientID uint   `gorm:"not null"`
	CityID   uint   `gorm:"not null"`
	Address  string `gorm:"not null"`
}
