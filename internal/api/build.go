package api

import (
	"fmt"
	"net/http"

	"github.com/fupas/commons/pkg/util"
	"github.com/labstack/echo/v4"
	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/analytics"
	"github.com/podops/podops/pkg/api"
	"github.com/podops/podops/pkg/auth"
	"github.com/podops/podops/pkg/backend"
	"google.golang.org/appengine"
)

// BuildEndpoint starts the build of the feed
func BuildEndpoint(c echo.Context) error {
	var req *a.Build = new(a.Build)

	if status, err := auth.Authorized(c, "ROLES"); err != nil {
		return api.ErrorResponse(c, status, err)
	}

	if err := c.Bind(req); err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	ctx := appengine.NewContext(c.Request())

	p, err := backend.GetProduction(ctx, req.GUID)
	if err != nil {
		return api.ErrorResponse(c, http.StatusNotFound, err)
	}
	if p == nil {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid guid '%s'", req.GUID))
	}

	// FIXME make this async, make validateOnly a flag
	if err := backend.Build(ctx, req.GUID, false); err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("error building feed '%s': %v", req.GUID, err))
	}

	// update the PRODUCTION record
	p.BuildDate = util.Timestamp()
	if err := backend.UpdateProduction(ctx, p); err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	resp := a.Build{
		GUID:         req.GUID,
		FeedURL:      fmt.Sprintf("%s/c/%s/feed.xml", a.DefaultCDNEndpoint, req.GUID),
		FeedAliasURL: fmt.Sprintf("%s/s/%s/feed.xml", a.DefaultPortalEndpoint, p.Name),
	}

	// track api access for billing etc
	analytics.TrackEvent(c.Request(), "api", "build", p.GUID, 1)

	return api.StandardResponse(c, http.StatusCreated, &resp)
}
