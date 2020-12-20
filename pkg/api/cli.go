package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/service/pkg/svc"

	"github.com/podops/podops/internal/cli"
)

// CreateNewShowEndpoint creates an new show
func CreateNewShowEndpoint(c *gin.Context) {
	var req cli.NewShowRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// create a show
	showName := strings.ToLower(strings.TrimSpace(req.Name))
	guid, _ := util.ShortUUID()

	// just send the ID back
	resp := cli.NewShowResponse{
		Name: showName,
		GUID: guid,
	}
	svc.StandardJSONResponse(c, &resp, nil)
}
