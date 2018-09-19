package models

import (
	"github.com/jinzhu/gorm"
)

type Client struct {
	gorm.Model
	Identification string `gorm:"not null; unique_index"`
	FirstName      string `gorm:"not null"`
	LastName       string `gorm:"not null"`
	Telephone1     string `gorm:"not null"`
	Telephone2     string
	Address        Address
}
