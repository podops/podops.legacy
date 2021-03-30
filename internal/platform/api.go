package platform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/appengine"
)

type (
	// StatusObject is used to report operation status and errors in an API request.
	// The struct can be used as a response object or be treated as an error object
	StatusObject struct {
		Status    int    `json:"status" binding:"required"`
		Message   string `json:"message" binding:"required"`
		RootError error  `json:"-"`
	}
)

// StandardResponse is the default way to respond to API requests
func StandardResponse(c echo.Context, status int, res interface{}) error {
	if res == nil {
		resp := StatusObject{
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
	var resp StatusObject

	// send the error to Google Error Reporting
	ReportError(err)

	if err == nil {
		resp = NewStatus(http.StatusInternalServerError, fmt.Sprintf("status: %d", status))
	} else {
		resp = NewErrorStatus(status, err)
	}
	return c.JSON(status, &resp)
}

// NewHttpContext creates a new context for appengine execution
// FIXME make this more pluggable or change outright
func NewHttpContext(c echo.Context) context.Context {
	return appengine.NewContext(c.Request())
}

// NewStatus initializes a new StatusObject
func NewStatus(s int, m string) StatusObject {
	return StatusObject{Status: s, Message: m}
}

// NewErrorStatus initializes a new StatusObject from an error
func NewErrorStatus(s int, e error) StatusObject {
	return StatusObject{Status: s, Message: e.Error(), RootError: e}
}

func (so *StatusObject) Error() string {
	return fmt.Sprintf("%s: %d", so.Message, so.Status)
}
