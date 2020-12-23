package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/platform/pkg/platform"

	t "github.com/podops/podops/internal/types"
)

// StandardResponse is the default way to respond to API requests
func StandardResponse(c *gin.Context, status int, res interface{}) {
	if res == nil {
		resp := t.StatusObject{
			Status:  status,
			Message: fmt.Sprintf("status: %d", status),
		}
		c.JSON(status, &resp)
	} else {
		c.JSON(status, res)
	}
}

// ErrorResponse responds with an ErrorObject
func ErrorResponse(c *gin.Context, status int, err error) {
	var resp t.StatusObject
	if err == nil {
		resp = t.StatusObject{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("status: %d", status), // keep it consistent with StandardResponse
		}
	} else {
		/*
			if ee, ok := err.(*t.StatusObject); ok {
				resp.Status = ee.Status
				resp.Message = ee.Message
			} else {
		*/
		resp = t.StatusObject{
			Status:  status,
			Message: err.Error(),
		}
		//}
	}

	c.JSON(status, &resp)
}

// HandleError is just a convenience method to avoid boiler-plate code
func HandleError(c *gin.Context, status int, e error) {
	platform.ReportError(e)
	ErrorResponse(c, status, e)
}
