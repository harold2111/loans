package main

import (
	"fmt"
	clientHtttpHandler "loans/client/handler/http"
	clientRepository "loans/client/repository/postgress"
	clientService "loans/client/service"
	"loans/config"
	"loans/errors"
	loanHtttpHandler "loans/loan/handler/http"
	loanRepository "loans/loan/repository/postgress"
	loanService "loans/loan/service"
	locationHtttpHandler "loans/location/handler/http"
	locationRepository "loans/location/repository/postgress"
	locationtService "loans/location/service"
	"loans/postgres"
	"loans/utils"
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

	clientRepository, _ := clientRepository.NewClientRepository(db)
	loanRepository, _ := loanRepository.NewLoanRepository(db)
	locationRepositoy, _ := locationRepository.NewLocationRepositoryy(db)

	clientService := clientService.NewClientService(clientRepository, locationRepositoy)
	loanService := loanService.NewLoanService(loanRepository, clientRepository)
	locationService := locationtService.NewLocationService(locationRepositoy)

	clientHtttpHandler.NewClientHttpHandler(echoContext, clientService)
	loanHtttpHandler.NewLoanHttpHandler(echoContext, loanService)
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
