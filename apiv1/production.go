package apiv1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/server"

	"github.com/podops/podops"
	"github.com/podops/podops/auth"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/errordef"
	lp "github.com/podops/podops/internal/platform"
)

// ProductionEndpoint creates an new show and does all the background setup
func ProductionEndpoint(c echo.Context) error {
	var req *podops.Production = new(podops.Production)
	ctx := platform.NewHttpContext(c.Request())

	if err := AuthorizeAccess(ctx, c, ScopeProductionWrite); err != nil {
		return server.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	err := c.Bind(req)
	if err != nil {
		return server.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// validate and normalize the name
	showName := strings.ToLower(strings.TrimSpace(req.Name))
	if !podops.ValidResourceName(showName) {
		return server.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(errordef.MsgInvalidParameter, showName))
	}
	// create a new production
	clientID, _ := auth.GetClientID(ctx, c.Request())
	p, err := backend.CreateProduction(ctx, showName, req.Title, req.Summary, clientID)
	if err != nil {
		return server.ErrorResponse(c, http.StatusBadRequest, err)
	}

	location := fmt.Sprintf("%s/show-%s.yaml", p.GUID, p.GUID)
	if err := backend.UpdateResource(ctx, p.Name, p.GUID, podops.ResourceShow, p.GUID, location); err != nil {
		return server.ErrorResponse(c, http.StatusBadRequest, err)
	}

	lp.TrackEvent(c.Request(), "api", "prod_create", p.GUID, 1)
	return server.StandardResponse(c, http.StatusCreated, p)
}

// ListProductionsEndpoint list all available shows
func ListProductionsEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c.Request())

	if err := AuthorizeAccess(ctx, c, ScopeProductionRead); err != nil {
		return server.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	clientID, _ := auth.GetClientID(ctx, c.Request())

	productions, err := backend.FindProductionsByOwner(ctx, clientID)
	if err != nil {
		return server.ErrorResponse(c, http.StatusBadRequest, err)
	}

	lp.TrackEvent(c.Request(), "api", "prod_list", clientID, 1)
	return server.StandardResponse(c, http.StatusOK, &podops.ProductionList{Productions: productions})
}
