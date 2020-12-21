package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// StandardJSONResponse is the default way to respond to API requests
func StandardJSONResponse(c *gin.Context, status int, res interface{}, err error) {
	if err == nil {
		if res == nil {
			c.JSON(status, gin.H{"status": fmt.Sprintf("%d", status)})
		} else {
			c.JSON(status, res)
		}
	} else {
		c.JSON(status, gin.H{"status": "error", "msg": err.Error()})
	}
}
