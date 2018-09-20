package models

import "time"

type Country struct {
	ID          uint         `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	DeletedAt   *time.Time   `sql:"index" json:"deletedAt"`
	Name        string       `gorm:"not null" json:"name"`
	Departments []Department `json:"departments"`
}
