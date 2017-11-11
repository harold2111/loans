package main

import (
	"fmt"
	"loans/client"
	"loans/config"
	"loans/errors"
	"loans/loan"
	"loans/postgres"
	"loans/utils"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func main() {

	time.Local = config.DefaultLocation()

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=loans sslmode=disable password=Nayarin1214")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	utils.InitValidator()
	postgres.MigrateModel(db)

	echoContext := echo.New()
	echoContext.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	clientRepository, _ := postgres.NewClientRepository(db)
	locationRepositoy, _ := postgres.NewLocationRepositoryy(db)

	clientService := client.NewService(clientRepository, locationRepositoy)

	client.SuscribeClientHandler(clientService, echoContext)
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
