package postgres

import (
	"loans/models"

	"github.com/jinzhu/gorm"
)

func MigrateModel(db *gorm.DB) {
	db.LogMode(true)

	db.DropTableIfExists(&models.Client{}, &models.Address{}, &models.City{}, &models.Department{},
		&models.Country{}, &models.Loan{}, &models.Bill{}, &models.BillMovement{}, &models.Payment{})

	db.CreateTable(&models.Client{}, &models.Address{}, &models.City{}, &models.Department{},
		&models.Country{}, &models.Loan{}, &models.Bill{}, &models.BillMovement{}, &models.Payment{})

	db.Model(&models.Client{}).Related(&models.Address{})

	cities := []models.City{
		{
			Name: "Barranquilla",
		},
	}
	departments := []models.Department{
		{
			Name:   "Atlántico",
			Cities: cities,
		},
	}
	country := models.Country{
		Name:        "Colombia",
		Departments: departments,
	}

	if error := db.Save(&country).Error; error != nil {
		panic(error)
	}

}
