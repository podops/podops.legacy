package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/appengine"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/platform"
)

// FIXME move this to package /internal/api

// StandardResponse is the default way to respond to API requests
func StandardResponse(c echo.Context, status int, res interface{}) error {
	if res == nil {
		resp := a.StatusObject{
			Status:  status,
			Message: fmt.Sprintf("status: %d", status),
		}
		return c.JSON(status, &resp)
	} else {
		return c.JSON(status, res)
	}
}

// ErrorResponse reports the error and responds with an ErrorObject
func ErrorResponse(c echo.Context, status int, err error) error {
	var resp a.StatusObject

	// send the error to Google Error Reporting
	platform.ReportError(err)

	if err == nil {
		resp = a.NewStatus(http.StatusInternalServerError, fmt.Sprintf("status: %d", status))
	} else {
		resp = a.NewErrorStatus(status, err)
	}
	return c.JSON(status, &resp)
}

// NewHttpContext creates a new context for appengine execution
// FIXME make this more pluggable or change outright
func NewHttpContext(c echo.Context) context.Context {
	return appengine.NewContext(c.Request())
}
