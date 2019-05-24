package models

import "time"

type Country struct {
	ID          uint         `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time    `json:"-"`
	UpdatedAt   time.Time    `json:"-"`
	DeletedAt   *time.Time   `sql:"index" json:"-"`
	Name        string       `gorm:"not null" json:"name"`
	Departments []Department `json:"departments"`
}
