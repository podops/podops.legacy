package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/platform/pkg/platform"

	"github.com/podops/podops/internal/errors"
)

const (
	// AdminNamespacePrefix namespace for internal admin endpoints
	AdminNamespacePrefix = "/_a"
	// NamespacePrefix namespace for the CLI. Should not be used directly.
	NamespacePrefix = "/a/v1"

	// All the API & CLI endpoint routes

	// AuthenticationRoute is used to create and verify a token
	AuthenticationRoute = "/token"
	// ProductionRoute route to ProductionEndpoint
	ProductionRoute = "/new"
	// ResourceRoute route to ResourceEndpoint
	ResourceRoute = "/update/:rsrc/:id"
)

// StandardResponse is the default way to respond to API requests
func StandardResponse(c *gin.Context, status int, res interface{}) {
	if res == nil {
		resp := errors.StatusObject{
			Status:  status,
			Message: fmt.Sprintf("Status %d", status),
		}
		c.JSON(status, &resp)
	} else {
		c.JSON(status, res)
	}
}

// ErrorResponse responds with an ErrorObject
func ErrorResponse(c *gin.Context, err error) {
	var resp errors.StatusObject
	if err == nil {
		resp = errors.StatusObject{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Status %d", http.StatusInternalServerError), // keep it consistent with StandardResponse
		}
	} else {
		if ee, ok := err.(*errors.StatusObject); ok {
			resp.Status = ee.Status
			resp.Message = ee.Message
		} else {
			resp = errors.StatusObject{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	}

	c.JSON(resp.Status, &resp)
}

// HandleError is just a convenience method to avoid boiler-plate code
func HandleError(c *gin.Context, e error) {
	platform.ReportError(e)
	ErrorResponse(c, e)
}
