package client

import (
	"github.com/jinzhu/gorm"
)

type Client struct {
	gorm.Model
	Identification string `gorm:"not null; unique_index"`
	FirstName      string `gorm:"not null"`
	LastName       string `gorm:"not null"`
	Telephone2     string
}

type Address struct {
	gorm.Model
	ClientID uint   `gorm:"not null"`
	CityID   uint   `gorm:"not null"`
	Address  string `gorm:"not null"`
}
