package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/auth"
	"google.golang.org/appengine"

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

// GetClientID extracts the ClientID from the token
func GetClientID(c *gin.Context) (string, error) {
	token := auth.GetBearerToken(c)
	if token == "" {
		return "", fmt.Errorf("production: missing token")
	}
	a, err := auth.FindAuthorization(appengine.NewContext(c.Request), token)
	if err != nil {
		return "", err
	}
	if a == nil {
		return "", fmt.Errorf("production: no authorization")
	}

	return a.ClientID, nil
}

// ExtractBodyAsString extracts a requests body, assuming it is a string
func ExtractBodyAsString(c *gin.Context) (string, error) {

	if c.Request.Body != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}
	return "", nil
}
