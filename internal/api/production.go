package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/platform"
	"github.com/podops/podops/pkg/api"
	"github.com/podops/podops/pkg/auth"
	"github.com/podops/podops/pkg/backend"
	"github.com/podops/podops/pkg/backend/models"
)

// ProductionEndpoint creates an new show and does all the background setup
func ProductionEndpoint(c echo.Context) error {
	var req *models.Production = new(models.Production)
	ctx := api.NewHttpContext(c)

	if err := AuthorizeAccess(ctx, c, scopeProductionWrite); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	err := c.Bind(req)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// validate and normalize the name
	showName := strings.ToLower(strings.TrimSpace(req.Name))
	if !a.ValidResourceName(showName) {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid name '%s'", showName))
	}
	// create a new production
	clientID, _ := auth.GetClientID(ctx, c.Request())
	p, err := backend.CreateProduction(ctx, showName, req.Title, req.Summary, clientID)
	if err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	location := fmt.Sprintf("%s/show-%s.yaml", p.GUID, p.GUID)
	if err := backend.UpdateResource(ctx, p.Name, p.GUID, a.ResourceShow, p.GUID, location); err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", "prod_create", p.GUID, 1)

	return api.StandardResponse(c, http.StatusCreated, p)
}

// ListProductionsEndpoint list all available shows
func ListProductionsEndpoint(c echo.Context) error {
	ctx := api.NewHttpContext(c)

	if err := AuthorizeAccess(ctx, c, scopeProductionRead); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	clientID, _ := auth.GetClientID(ctx, c.Request())

	productions, err := backend.FindProductionsByOwner(ctx, clientID)
	if err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", "prod_list", clientID, 1)

	return api.StandardResponse(c, http.StatusOK, &models.ProductionList{Productions: productions})
}
