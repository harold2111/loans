package postgres

import (
	clientDomain "github.com/harold2111/loans/client/domain"
	loanDomain "github.com/harold2111/loans/loan/domain"
	locationDomain "github.com/harold2111/loans/location/domain"

	"github.com/jinzhu/gorm"
)

func MigrateModel(db *gorm.DB) {
	db.LogMode(true)

	db.DropTableIfExists(
		//Client module
		&clientDomain.Client{},
		&clientDomain.Address{},
		//Location module
		&locationDomain.City{},
		&locationDomain.Department{},
		&locationDomain.Country{},
		//Loan module
		&loanDomain.Loan{},
		&loanDomain.Payment{},
		&loanDomain.Period{},
		&loanDomain.PeriodPayment{},
		&loanDomain.DefaultPeriod{},
		&loanDomain.DefaultPayment{},
	)
	db.CreateTable( //Client module
		&clientDomain.Client{},
		&clientDomain.Address{},
		//Location module
		&locationDomain.City{},
		&locationDomain.Department{},
		&locationDomain.Country{},
		//Loan module
		&loanDomain.Loan{},
		&loanDomain.Payment{},
		&loanDomain.Period{},
		&loanDomain.PeriodPayment{},
		&loanDomain.DefaultPeriod{},
		&loanDomain.DefaultPayment{},
	)

	db.Model(&clientDomain.Client{}).Related(&clientDomain.Address{})

	atlanticoCities := []locationDomain.City{
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
	antioquiaCities := []locationDomain.City{
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
	departments := []locationDomain.Department{
		{
			Name:   "Atlántico",
			Cities: atlanticoCities,
		},
		{
			Name:   "Antioquia",
			Cities: antioquiaCities,
		},
	}
	country := locationDomain.Country{
		Name:        "Colombia",
		Departments: departments,
	}

	if error := db.Save(&country).Error; error != nil {
		panic(error)
	}

}
