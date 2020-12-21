package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/platform/pkg/platform"

	"github.com/podops/podops/internal/cli"
	"github.com/podops/podops/internal/production"
)

// NewShowEndpoint creates an new show and does all the background setup
func NewShowEndpoint(c *gin.Context) {
	var req cli.NewShowRequest

	err := c.BindJSON(&req)
	if err != nil {
		platform.ReportError(err)
		c.Status(http.StatusBadRequest)
		return
	}

	// create a show
	showName := strings.ToLower(strings.TrimSpace(req.Name)) // FIXME: verify && cleanup the name. Should follow Domain name conventions.

	p, status, err := production.CreateProduction(appengine.NewContext(c.Request), showName, req.Title, req.Summary)
	if err != nil {
		platform.ReportError(err)
		StandardJSONResponse(c, status, nil, err)
		return
	}

	// just send the GUID and canonical name back
	resp := cli.NewShowResponse{
		Name: showName,
		GUID: p.GUID,
	}
	StandardJSONResponse(c, status, &resp, nil)
}
