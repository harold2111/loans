package models

import "github.com/jinzhu/gorm"

type City struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Department   Department
	DepartmentID uint `gorm:"not null"`
}
