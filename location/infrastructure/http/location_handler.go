package http

import (
	locationApplication "loans/location/application"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
)

type HttplocationHandler struct {
	locationService locationApplication.LocationService
}

func NewLocationHttpHandler(e *echo.Echo, locationService locationApplication.LocationService) {
	handler := &HttplocationHandler{
		locationService: locationService,
	}
	e.GET("/api/departments", handler.handleFindAllDepartments)
	e.GET("/api/cities", handler.handleFindCitiesByDepartmentID)
}

func (handler *HttplocationHandler) handleFindAllDepartments(c echo.Context) error {
	locationService := handler.locationService
	departments, error := locationService.FindAllDepartments()
	if error != nil {
		return error
	}
	response := new([]locationApplication.DepartmentResponse)
	if error := copier.Copy(&response, &departments); error != nil {
		return error
	}
	return c.JSON(http.StatusOK, response)
}

func (handler *HttplocationHandler) handleFindCitiesByDepartmentID(c echo.Context) error {
	locationService := handler.locationService
	deparmentID, _ := strconv.Atoi(c.QueryParam("departmentID"))
	cities, error := locationService.FindCitiesByDepartmentID(uint(deparmentID))
	if error != nil {
		return error
	}
	response := new([]locationApplication.CityResponse)
	if error := copier.Copy(&response, &cities); error != nil {
		return error
	}
	return c.JSON(http.StatusOK, response)
}
