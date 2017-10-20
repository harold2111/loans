package errors

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	var (
		code             = http.StatusInternalServerError
		msg  interface{} = err.Error()
	)
	switch errorType := err.(type) {
	case *RecordNotFound:
		code = http.StatusNotFound
	case *GracefulError:
		code = http.StatusBadRequest
	case *ValidationError:
		code = http.StatusBadRequest
	case *echo.HTTPError:
		code = errorType.Code
		msg = errorType.Message
		if errorType.Inner != nil {
			msg = fmt.Sprintf("%v, %v", errorType, errorType.Inner)
		}
	}

	if _, ok := msg.(string); ok {
		msg = echo.Map{"message": msg}
	} else {
		msg = echo.Map{"message": http.StatusText(code)}
	}

	c.Logger().Error(err)

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
