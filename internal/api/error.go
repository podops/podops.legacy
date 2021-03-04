package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/pkg/observer"
)

// ErrorResponse reports the error and responds with an ErrorObject
func ErrorResponse(c echo.Context, status int, err error) error {
	var resp a.StatusObject

	// send the error to Google Error Reporting
	observer.ReportError(err)

	if err == nil {
		resp = a.NewStatus(http.StatusInternalServerError, fmt.Sprintf("status: %d", status))
	} else {
		resp = a.NewErrorStatus(status, err)
	}
	return c.JSON(status, &resp)
}
