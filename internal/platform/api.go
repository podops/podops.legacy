package platform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podops/podops"
	"google.golang.org/appengine"
)

// StandardResponse is the default way to respond to API requests
func StandardResponse(c echo.Context, status int, res interface{}) error {
	if res == nil {
		resp := podops.StatusObject{
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
	var resp podops.StatusObject

	// send the error to Google Error Reporting
	ReportError(err)

	if err == nil {
		resp = podops.NewStatus(http.StatusInternalServerError, fmt.Sprintf("status: %d", status))
	} else {
		resp = podops.NewErrorStatus(status, err)
	}
	return c.JSON(status, &resp)
}

// NewHttpContext creates a new context for appengine execution
// FIXME make this more pluggable or change outright
func NewHttpContext(c echo.Context) context.Context {
	return appengine.NewContext(c.Request())
}
