package migration

import (
	"loans/config"
	"loans/models"

	"github.com/jinzhu/gorm"
)

func MigrateModel(db *gorm.DB) {
	db.LogMode(true)
	db.DropTableIfExists(&models.Client{}, &models.Address{}, &models.City{}, &models.Department{}, &models.Country{}, &models.Product{}, &models.Loan{})
	db.CreateTable(&models.Address{}, &models.Client{}, &models.City{}, &models.Department{}, &models.Country{}, &models.Product{}, &models.Loan{})

	db.Model(&models.Client{}).Related(&models.Address{})

	country := models.Country{Name: "Colombia"}
	department := models.Department{Name: "Atlántico", Country: country}
	city := models.City{Name: "Barranquilla", Department: department}

	if error := config.DB.Save(&city).Error; error != nil {
		panic(error)
	}

}
