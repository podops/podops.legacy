package apiv1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fupas/commons/pkg/env"
	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/api"
	"github.com/txsvc/platform/pkg/tasks"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/feed"
	"github.com/podops/podops/internal/errordef"
	lp "github.com/podops/podops/internal/platform"
)

var (
	// full canonical route
	syncTaskEndpoint string = podops.DefaultCDNEndpoint + "/_w/sync"
)

// BuildFeedEndpoint starts the build of the feed
func BuildFeedEndpoint(c echo.Context) error {
	var req *podops.BuildRequest = new(podops.BuildRequest) // FIXME change this
	ctx := platform.NewHttpContext(c.Request())

	if err := c.Bind(req); err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if err := AuthorizeAccessProduction(ctx, c, ScopeProductionBuild, req.GUID); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	validateOnly := false
	if strings.ToLower(c.QueryParam("v")) == "true" {
		validateOnly = true
	}

	p, err := backend.GetProduction(ctx, req.GUID)
	if err != nil {
		return api.ErrorResponse(c, http.StatusNotFound, err)
	}
	if p == nil {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(errordef.MsgInvalidGUID, req.GUID))
	}

	if err := feed.Build(ctx, req.GUID, validateOnly); err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	if !validateOnly {
		// dispatch a request for background sync
		ir := podops.SyncRequest{
			GUID:   req.GUID,
			Source: "feed.xml",
		}

		task := tasks.HttpTask{
			Method:  tasks.HttpMethodPost,
			Request: syncTaskEndpoint,
			Token:   env.GetString("PODOPS_API_KEY", ""),
			Payload: &ir,
		}
		err := platform.NewTask(task)
		if err != nil {
			return err
		}
	}

	// track api access for billing etc
	lp.TrackEvent(c.Request(), "api", "build", p.GUID, 1)

	resp := podops.BuildRequest{
		GUID:         req.GUID,
		FeedURL:      fmt.Sprintf("%s/%s/feed.xml", podops.DefaultStorageEndpoint, req.GUID),
		FeedAliasURL: fmt.Sprintf("%s/s/%s/feed.xml", podops.DefaultEndpoint, p.Name),
	}

	return api.StandardResponse(c, http.StatusCreated, &resp)
}
