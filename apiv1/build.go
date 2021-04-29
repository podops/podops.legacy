package apiv1

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform/pkg/env"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/feed"
	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/platform"
)

var (
	// full canonical route
	syncTaskEndpoint string = podops.DefaultCDNEndpoint + "/_w/sync"
)

// BuildFeedEndpoint starts the build of the feed
func BuildFeedEndpoint(c echo.Context) error {
	var req *podops.BuildRequest = new(podops.BuildRequest) // FIXME change this
	ctx := platform.NewHttpContext(c)

	if err := c.Bind(req); err != nil {
		return platform.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if err := AuthorizeAccessProduction(ctx, c, ScopeProductionBuild, req.GUID); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	validateOnly := false
	if strings.ToLower(c.QueryParam("v")) == "true" {
		validateOnly = true
	}

	p, err := backend.GetProduction(ctx, req.GUID)
	if err != nil {
		return platform.ErrorResponse(c, http.StatusNotFound, err)
	}
	if p == nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(errordef.MsgInvalidGUID, req.GUID))
	}

	if err := feed.Build(ctx, req.GUID, validateOnly); err != nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, err)
	}

	if !validateOnly {
		// dispatch a request for background sync
		ir := podops.SyncRequest{
			GUID:   req.GUID,
			Source: "feed.xml",
		}
		_, err = platform.CreateHttpTask(ctx, tasks.HttpMethod_POST, syncTaskEndpoint, env.GetString("PODOPS_API_KEY", ""), &ir)
		if err != nil {
			return err
		}
	}

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", "build", p.GUID, 1)

	resp := podops.BuildRequest{
		GUID:         req.GUID,
		FeedURL:      fmt.Sprintf("%s/%s/feed.xml", podops.DefaultStorageEndpoint, req.GUID),
		FeedAliasURL: fmt.Sprintf("%s/s/%s/feed.xml", podops.DefaultEndpoint, p.Name),
	}

	return platform.StandardResponse(c, http.StatusCreated, &resp)
}
