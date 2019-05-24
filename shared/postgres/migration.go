package postgres

import (
	"loans/shared/models"

	"github.com/jinzhu/gorm"
)

func MigrateModel(db *gorm.DB) {
	db.LogMode(true)

	db.DropTableIfExists(&models.Client{}, &models.Address{}, &models.City{}, &models.Department{},
		&models.Country{}, &models.Loan{}, &models.Bill{}, &models.BillMovement{}, &models.Payment{})

	db.CreateTable(&models.Client{}, &models.Address{}, &models.City{}, &models.Department{},
		&models.Country{}, &models.Loan{}, &models.Bill{}, &models.BillMovement{}, &models.Payment{})

	db.Model(&models.Client{}).Related(&models.Address{})

	atlanticoCities := []models.City{
		{
			Name: "Barranquilla",
		},
		{
			Name: "Soledad",
		},
		{
			Name: "Pto Colombia",
		},
	}
	antioquiaCities := []models.City{
		{
			Name: "Medellin",
		},
		{
			Name: "Envigado",
		},
		{
			Name: "Sabaneta",
		},
	}
	departments := []models.Department{
		{
			Name:   "Atlántico",
			Cities: atlanticoCities,
		},
		{
			Name:   "Antioquia",
			Cities: antioquiaCities,
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
