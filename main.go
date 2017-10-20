package main

import (
	"loans/config"
	"loans/controllers"
	"loans/errors"
	"loans/migration"
	"loans/validators"

	"github.com/labstack/echo"
)

func main() {
	validators.InitValidator()
	config.InitDB("host=localhost user=postgres dbname=loans sslmode=disable password=Nayarin1214")
	migration.MigrateModel(config.DB)

	echoContext := echo.New()
	echoContext.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	echoContext.POST("/api/clients", controllers.CreateClient)
	echoContext.PUT("/api/clients/:id", controllers.UpdateClient)

	echoContext.Logger.Fatal(echoContext.Start(":1323"))

}
