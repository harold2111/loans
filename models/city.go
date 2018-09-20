package models

import "time"

type City struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `sql:"index" json:"deletedAt"`
	Name         string     `gorm:"not null" json:"name"`
	DepartmentID uint       `gorm:"not null" json:"departmentID"`
}
