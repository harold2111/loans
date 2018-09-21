package models

import "time"

type Client struct {
	ID             uint       `gorm:"primary_key" json:"id"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `sql:"index" json:"deletedAt"`
	Identification string     `gorm:"not null; unique_index"`
	FirstName      string     `gorm:"not null" json:"firstName" validate:"required"`
	LastName       string     `gorm:"not null" json:"lastName" validate:"required"`
	Telephone1     string     `gorm:"not null" json:"telephone1" validate:"required"`
	Telephone2     string     `json:"telephone2"`
	Addresses      []Address  `validate:"required,dive,required" json:"addresses,omitempty" `
}
