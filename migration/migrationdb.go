package migration

import (
	"loans/address"
	"loans/client"
	"loans/config"
	"loans/loan"

	"github.com/jinzhu/gorm"
)

func MigrateModel(db *gorm.DB) {
	db.LogMode(true)

	db.DropTableIfExists(&client.Client{}, &address.Address{}, &address.City{}, &address.Department{},
		&address.Country{}, &loan.Loan{}, &loan.Bill{}, &loan.BillMovement{}, &loan.Payment{})

	db.CreateTable(&client.Client{}, &address.Address{}, &address.City{}, &address.Department{},
		&address.Country{}, &loan.Loan{}, &loan.Bill{}, &loan.BillMovement{}, &loan.Payment{})

	db.Model(&client.Client{}).Related(&address.Address{})

	country := address.Country{Name: "Colombia"}
	department := address.Department{Name: "Atlántico", Country: country}
	city := address.City{Name: "Barranquilla", Department: department}

	if error := config.DB.Save(&city).Error; error != nil {
		panic(error)
	}

}
