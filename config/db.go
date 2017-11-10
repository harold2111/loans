package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

var DB *gorm.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = gorm.Open("postgres", dataSourceName)
	if err != nil {
		panic("failed to connect database")
	}
}

func IsUniqueConstraintError(err error, constraintName string) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" && pqErr.Constraint == constraintName
	}
	return false
}
