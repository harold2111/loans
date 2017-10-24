package main

import (
	"fmt"
	"loans/config"
	"loans/controllers"
	"loans/errors"
	"loans/migration"
	"loans/validators"
	"time"

	"github.com/labstack/echo"
)

func main() {

	validators.InitValidator()
	config.InitDB("host=localhost user=postgres dbname=loans sslmode=disable password=Nayarin1214")
	migration.MigrateModel(config.DB)
	startPaymentJob()

	echoContext := echo.New()
	echoContext.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	echoContext.POST("/api/clients", controllers.CreateClient)
	echoContext.PUT("/api/clients/:id", controllers.UpdateClient)

	echoContext.POST("/api/loans", controllers.CreateLoan)

	echoContext.Logger.Fatal(echoContext.Start(":1323"))

}

func startPaymentJob() {
	ticker := time.NewTicker(1 * time.Second)
	doneChan := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println(time.Now())
			case <-doneChan:
				ticker.Stop()
				fmt.Println("Stopped the ticker!")
				return
			}
		}
	}()
	//close(quit)
}
