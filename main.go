package main

import (
	"fmt"
	"loans/config"
	"loans/controllers"
	"loans/errors"
	"loans/migration"
	"loans/validators"

	"github.com/shopspring/decimal"

	"github.com/labstack/echo"
)

func main() {

	validators.InitValidator()
	config.InitDB("host=localhost user=postgres dbname=loans sslmode=disable password=Nayarin1214")
	migration.MigrateModel(config.DB)

	var total float64 = 0
	fmt.Println(total)
	total += 5.6
	total += 5.8
	fmt.Println(total)

	var n1 float64 = 10
	fmt.Println(n1)
	var n2 float64 = 3
	var n3 float64 = n1 / n2
	fmt.Println(n3)

	nu1, _ := decimal.NewFromString("10")
	nu2, _ := decimal.NewFromString("3")

	nu3 := nu1.Div(nu2)

	fmt.Println(nu3.String())

	echoContext := echo.New()
	echoContext.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	echoContext.POST("/api/clients", controllers.CreateClient)
	echoContext.PUT("/api/clients/:id", controllers.UpdateClient)

	echoContext.POST("/api/loans", controllers.CreateLoan)

	echoContext.Logger.Fatal(echoContext.Start(":1323"))

}
