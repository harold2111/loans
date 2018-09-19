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

	country := models.Country{Name: "Colombia"}
	department := models.Department{Name: "Atlántico", Country: country}
	city := models.City{Name: "Barranquilla", Department: department}

	if error := db.Save(&city).Error; error != nil {
		panic(error)
	}

}
