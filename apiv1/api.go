package apiv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fupas/observer/pkg/observer"
	"github.com/gin-gonic/gin"
)

const (
	// Version specifies the verion of the API and its structs
	Version = "v1"

	// MajorVersion of the API
	MajorVersion = 1
	// MinorVersion of the API
	MinorVersion = 1
	// FixVersion of the API
	FixVersion = 2
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

var (
	// ErrNotAuthorized indicates that the API call is not authorized
	ErrNotAuthorized = errors.New("api: not authorized")
	// ErrNoToken indicates that no bearer token was provided
	ErrNoToken = errors.New("api: no token provided")

	// ErrInvalidParameters indicates that parameters used in an API call are not valid
	ErrInvalidParameters = errors.New("api: invalid parameters")
	// ErrValidationFailed indicates that a resource validation failed
	ErrValidationFailed = errors.New("api: validation failed")

	// ErrNoSuchProduction indicates that the production does not exist
	ErrNoSuchProduction = errors.New("api: production doesn't exist")
	// ErrNoSuchResource indicates that the resource does not exist
	ErrNoSuchResource = errors.New("api: resource doesn't exist")
	// ErrNoSuchAsset indicates that the asset does not exist
	ErrNoSuchAsset = errors.New("api: asset doesn't exist")
	// ErrBuildFailed indicates that the feed build failed
	ErrBuildFailed = errors.New("api: build failed")

	// ErrInternalError indicates that an unspecified internal error happened
	ErrInternalError = errors.New("api: internal error")

	// VersionString is the canonical API description
	VersionString string = fmt.Sprintf("%s.%s.%s", MajorVersion, MinorVersion, FixVersion)
)

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

// StandardResponse is the default way to respond to API requests
func StandardResponse(c *gin.Context, status int, res interface{}) {
	if res == nil {
		resp := StatusObject{
			Status:  status,
			Message: fmt.Sprintf("status: %d", status),
		}
		c.JSON(status, &resp)
	} else {
		c.JSON(status, res)
	}
}

// ErrorResponse reports the error and responds with an ErrorObject
func ErrorResponse(c *gin.Context, status int, err error) {
	var resp StatusObject

	// send the error to Google Error Reporting
	observer.ReportError(err)

	if err == nil {
		resp = NewStatus(http.StatusInternalServerError, fmt.Sprintf("status: %d", status))
	} else {
		resp = NewErrorStatus(status, err)
	}
	c.JSON(status, &resp)
}

// VersionEndpoint returns the current API version
func VersionEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": VersionString, "major": MajorVersion, "minor": MinorVersion, "fix": FixVersion, "namespace": Version})
}
