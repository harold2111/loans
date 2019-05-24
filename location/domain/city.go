package domain

import "time"

type City struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    *time.Time `sql:"index" json:"-"`
	Name         string     `gorm:"not null" json:"name"`
	DepartmentID uint       `gorm:"not null" json:"departmentID"`
}
