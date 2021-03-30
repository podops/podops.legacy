package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	a "github.com/podops/podops"
	"github.com/podops/podops/internal/platform"
	"github.com/podops/podops/pkg/auth"
	"github.com/podops/podops/pkg/backend"
)

// ProductionEndpoint creates an new show and does all the background setup
func ProductionEndpoint(c echo.Context) error {
	var req *a.Production = new(a.Production)
	ctx := platform.NewHttpContext(c)

	if err := AuthorizeAccess(ctx, c, scopeProductionWrite); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	err := c.Bind(req)
	if err != nil {
		return platform.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// validate and normalize the name
	showName := strings.ToLower(strings.TrimSpace(req.Name))
	if !a.ValidResourceName(showName) {
		return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid name '%s'", showName))
	}
	// create a new production
	clientID, _ := auth.GetClientID(ctx, c.Request())
	p, err := backend.CreateProduction(ctx, showName, req.Title, req.Summary, clientID)
	if err != nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, err)
	}

	location := fmt.Sprintf("%s/show-%s.yaml", p.GUID, p.GUID)
	if err := backend.UpdateResource(ctx, p.Name, p.GUID, a.ResourceShow, p.GUID, location); err != nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", "prod_create", p.GUID, 1)

	return platform.StandardResponse(c, http.StatusCreated, p)
}

// ListProductionsEndpoint list all available shows
func ListProductionsEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c)

	if err := AuthorizeAccess(ctx, c, scopeProductionRead); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	clientID, _ := auth.GetClientID(ctx, c.Request())

	productions, err := backend.FindProductionsByOwner(ctx, clientID)
	if err != nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", "prod_list", clientID, 1)

	return platform.StandardResponse(c, http.StatusOK, &a.ProductionList{Productions: productions})
}
