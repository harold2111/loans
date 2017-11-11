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
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	loanRepository, _ := postgres.NewLoanRepository(db)
	locationRepositoy, _ := postgres.NewLocationRepositoryy(db)

	clientService := client.NewService(clientRepository, locationRepositoy)
	loanService := loan.NewService(loanRepository, clientRepository)

	client.SuscribeClientHandler(clientService, echoContext)
	loan.SuscribeLoanHandler(loanService, echoContext)

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
