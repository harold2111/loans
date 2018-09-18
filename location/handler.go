package location

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
)

func SuscribeLocationHandler(s Service, e *echo.Echo) {
	e.GET("/api/departments", func(c echo.Context) error {
		return handleFindAllDepartments(s, c)
	})
	e.GET("/api/cities", func(c echo.Context) error {
		return handleFindCitiesByDepartmentID(s, c)
	})
}

func handleFindAllDepartments(s Service, c echo.Context) error {
	departments, error := s.FindAllDepartments()
	if error != nil {
		return error
	}
	response := new([]DepartmentResponse)
	if error := copier.Copy(&response, &departments); error != nil {
		return error
	}
	return c.JSON(http.StatusOK, response)
}

func handleFindCitiesByDepartmentID(s Service, c echo.Context) error {
	deparmentID, _ := strconv.Atoi(c.QueryParam("departmentID"))
	cities, error := s.FindCitiesByDepartmentID(uint(deparmentID))
	if error != nil {
		return error
	}
	response := new([]CityResponse)
	if error := copier.Copy(&response, &cities); error != nil {
		return error
	}
	return c.JSON(http.StatusOK, response)
}
