package models

import "time"

type Client struct {
	ID             uint       `gorm:"primary_key" json:"id"`
	CreatedAt      time.Time  `json:"-"`
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `sql:"index" json:"-"`
	Identification string     `gorm:"not null; unique_index"`
	FirstName      string     `gorm:"not null" json:"firstName" validate:"required"`
	LastName       string     `gorm:"not null" json:"lastName" validate:"required"`
	Telephone1     string     `gorm:"not null" json:"telephone1" validate:"required"`
	Email          string     `json:"email" validate:"required"`
	Addresses      []Address  `validate:"dive" json:"addresses,omitempty" `
}
