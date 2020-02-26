package main

import (
	"fmt"
	"time"

	clientApplication "github.com/harold2111/loans/client/application"
	clientHttpHandler "github.com/harold2111/loans/client/infrastructure/http"
	clientPostgresRepository "github.com/harold2111/loans/client/infrastructure/postgress"
	loanApplication "github.com/harold2111/loans/loan/application"
	loanHttpHandler "github.com/harold2111/loans/loan/infrastructure/http"
	loanPostgresRepository "github.com/harold2111/loans/loan/infrastructure/postgress"
	locationApplication "github.com/harold2111/loans/location/application"
	locationHttpHandler "github.com/harold2111/loans/location/infrastructure/http"
	locationPostgresRepository "github.com/harold2111/loans/location/infrastructure/postgress"
	"github.com/harold2111/loans/shared/config"
	"github.com/harold2111/loans/shared/errors"
	"github.com/harold2111/loans/shared/postgres"
	"github.com/harold2111/loans/shared/utils"

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

	clientRepository, _ := clientPostgresRepository.NewClientRepository(db)
	loanRepository, _ := loanPostgresRepository.NewLoanRepository(db)
	locationRepository, _ := locationPostgresRepository.NewLocationRepository(db)

	clientService := clientApplication.NewClientService(clientRepository, locationRepository)
	loanService := loanApplication.NewLoanService(loanRepository, clientRepository)
	locationService := locationApplication.NewLocationService(locationRepository)

	clientHttpHandler.NewClientHttpHandler(echoContext, clientService)
	loanHttpHandler.NewLoanHttpHandler(echoContext, loanService)
	locationHttpHandler.NewLocationHttpHandler(echoContext, locationService)

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
