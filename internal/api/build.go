package api

import (
	"fmt"
	"net/http"

	"github.com/fupas/commons/pkg/util"
	"github.com/labstack/echo/v4"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/platform"
	"github.com/podops/podops/pkg/api"
	"github.com/podops/podops/pkg/backend"
)

// BuildFeedEndpoint starts the build of the feed
func BuildFeedEndpoint(c echo.Context) error {
	var req *a.BuildRequest = new(a.BuildRequest)
	ctx := api.NewHttpContext(c)

	if err := c.Bind(req); err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if err := AuthorizeAccessProduction(ctx, c, scopeProductionBuild, req.GUID); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	p, err := backend.GetProduction(ctx, req.GUID)
	if err != nil {
		return api.ErrorResponse(c, http.StatusNotFound, err)
	}
	if p == nil {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid guid '%s'", req.GUID))
	}

	// FIXME make this async, make validateOnly a flag
	if err := backend.BuildFeed(ctx, req.GUID, false); err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("error building feed '%s': %v", req.GUID, err))
	}

	// update the PRODUCTION record
	p.BuildDate = util.Timestamp()
	if err := backend.UpdateProduction(ctx, p); err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	resp := a.BuildRequest{
		GUID:         req.GUID,
		FeedURL:      fmt.Sprintf("%s/c/%s/feed.xml", a.DefaultCDNEndpoint, req.GUID),
		FeedAliasURL: fmt.Sprintf("%s/s/%s/feed.xml", a.DefaultPortalEndpoint, p.Name),
	}

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", "build", p.GUID, 1)

	return api.StandardResponse(c, http.StatusCreated, &resp)
}
