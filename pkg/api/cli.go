package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/service/pkg/svc"

	"github.com/podops/podops/internal/cli"
)

// CreateNewShowEndpoint creates an new show
func CreateNewShowEndpoint(c *gin.Context) {
	var req cli.CLINewShowRequest
	var resp cli.CLINewShowResponse

	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resp.GUID, _ = util.ShortUUID()

	svc.StandardJSONResponse(c, &resp, nil)
}
