package models

import "time"

type Address struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `sql:"index" json:"deletedAt"`
	Address      string     `gorm:"not null" json:"address" validate:"required"`
	ClientID     uint       `gorm:"not null" json:"clientID"`
	DepartmentID uint       `gorm:"not null" json:"departmentID"`
	CityID       uint       `gorm:"not null" json:"cityID" validate:"required"`
}
