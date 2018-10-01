package models

import "time"

type Address struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    *time.Time `sql:"index" json:"-"`
	StretAddress string     `gorm:"not null" json:"stretAddress" validate:"required"`
	ClientID     uint       `gorm:"not null" json:"clientID"`
	DepartmentID uint       `gorm:"not null" json:"departmentID" validate:"required"`
	CityID       uint       `gorm:"not null" json:"cityID" validate:"required"`
}
