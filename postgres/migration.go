package postgres

import (
	"loans/client"
	"loans/loan"
	"loans/location"

	"github.com/jinzhu/gorm"
)

func MigrateModel(db *gorm.DB) {
	db.LogMode(true)

	db.DropTableIfExists(&client.Client{}, &client.Address{}, &location.City{}, &location.Department{},
		&location.Country{}, &loan.Loan{}, &loan.Bill{}, &loan.BillMovement{}, &loan.Payment{})

	db.CreateTable(&client.Client{}, &client.Address{}, &location.City{}, &location.Department{},
		&location.Country{}, &loan.Loan{}, &loan.Bill{}, &loan.BillMovement{}, &loan.Payment{})

	db.Model(&client.Client{}).Related(&client.Address{})

	country := location.Country{Name: "Colombia"}
	department := location.Department{Name: "Atlántico", Country: country}
	city := location.City{Name: "Barranquilla", Department: department}

	if error := db.Save(&city).Error; error != nil {
		panic(error)
	}

}
