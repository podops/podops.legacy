package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/service/pkg/svc"

	"github.com/podops/podops/internal/cli"
)

// NewShowEndpoint creates an new show and does all the background setup
func NewShowEndpoint(c *gin.Context) {
	var req cli.NewShowRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// create a show
	showName := strings.ToLower(strings.TrimSpace(req.Name)) // FIXME: verify && cleanup the name. Should follow Domain name conventions.
	guid, _ := util.ShortUUID()

	// # FIXME create the actual setup

	// just send the ID back
	resp := cli.NewShowResponse{
		Name: showName,
		GUID: guid,
	}
	svc.StandardJSONResponse(c, &resp, nil)
}
