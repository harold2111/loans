package location

import "github.com/jinzhu/gorm"

type Department struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Country   Country
	CountryID int `gorm:"not null"`
}

type City struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Department   Department
	DepartmentID uint `gorm:"not null"`
}

type Country struct {
	gorm.Model
	Name string `gorm:"not null"`
}
