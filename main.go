package main

import (
	"fmt"
	clientApplication "loans/client/application"
	clientHttpHandler "loans/client/infrastructure/http"
	clientPostgressRepository "loans/client/infrastructure/postgress"
	loanApplication "loans/loan/application"
	loantHttpHandler "loans/loan/infrastructure/http"
	loanPostgressRepository "loans/loan/infrastructure/postgress"
	locationApplication "loans/location/application"
	locationHtttpHandler "loans/location/infrastructure/http"
	locationPostgressRepository "loans/location/infrastructure/postgress"
	"loans/shared/config"
	"loans/shared/errors"
	"loans/shared/postgres"
	"loans/shared/utils"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	echoContext.Use(middleware.CORS())
	echoContext.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	clientRepository, _ := clientPostgressRepository.NewClientRepository(db)
	loanRepository, _ := loanPostgressRepository.NewLoanRepository(db)
	locationRepositoy, _ := locationPostgressRepository.NewLocationRepositoryy(db)

	clientService := clientApplication.NewClientService(clientRepository, locationRepositoy)
	loanService := loanApplication.NewLoanService(loanRepository, clientRepository)
	locationService := locationApplication.NewLocationService(locationRepositoy)

	clientHttpHandler.NewClientHttpHandler(echoContext, clientService)
	loantHttpHandler.NewLoanHttpHandler(echoContext, loanService)
	locationHtttpHandler.NewLocationHttpHandler(echoContext, locationService)

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
