package main

import (
	"fmt"
	"loans/client"
	"loans/config"
	"loans/errors"
	"loans/loan"
	"loans/migration"
	"loans/utils"
	"time"

	"github.com/labstack/echo"
)

func main() {

	time.Local = config.DefaultLocation()

	utils.InitValidator()
	config.InitDB("host=localhost user=postgres dbname=loans sslmode=disable password=Nayarin1214")
	migration.MigrateModel(config.DB)
	//startPaymentJob()

	echoContext := echo.New()
	echoContext.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	echoContext.POST("/api/clients", client.CreateClient)
	echoContext.PUT("/api/clients/:id", client.UpdateClient)

	echoContext.POST("/api/loans", loan.CreateLoan)

	echoContext.POST("/api/loans/payments", loan.PayLoan)

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
