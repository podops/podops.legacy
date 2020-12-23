package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/podops/podops/internal/production"
	t "github.com/podops/podops/internal/types"
)

// ProductionEndpoint creates an new show and does all the background setup
func ProductionEndpoint(c *gin.Context) {
	var req t.ProductionRequest

	err := c.BindJSON(&req)
	if err != nil {
		HandleError(c, err)
		return
	}

	// create a show
	// FIXME: verify && cleanup the name. Should follow Domain name conventions.
	showName := strings.ToLower(strings.TrimSpace(req.Name))
	p, err := production.CreateProduction(appengine.NewContext(c.Request), showName, req.Title, req.Summary)
	if err != nil {
		HandleError(c, err)
		return
	}

	// send the GUID and canonical name back
	resp := t.ProductionResponse{
		Name: p.Name,
		GUID: p.GUID,
	}
	StandardResponse(c, http.StatusCreated, &resp)
}
