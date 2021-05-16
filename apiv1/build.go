package apiv1

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/pkg/api"
	"github.com/txsvc/platform/v2/pkg/apis/provider"
	"github.com/txsvc/platform/v2/pkg/env"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/feed"
	"github.com/podops/podops/internal/messagedef"
)

var (
	// full canonical route
	syncTaskEndpoint string = podops.DefaultCDNEndpoint + "/_w/sync"

	tp provider.HttpTaskProvider
)

// implements lazy loading to give other parts of the code time to initialize the platform
// before a first call to the authentication provider is made. This is why init() would not work.

func background() provider.HttpTaskProvider {
	if tp != nil {
		return tp
	}
	p, ok := platform.Provider(provider.TypeTask)
	if !ok {
		err := fmt.Errorf(platform.MsgMissingProvider, provider.TypeTask.String())
		platform.ReportError(err)
		log.Fatal(err) // this halts the process but there is no point because it would just crash later anyways
	}
	tp = p.(provider.HttpTaskProvider)

	return tp
}

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
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(messagedef.MsgResourceInvalidGUID, req.GUID))
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

		task := provider.HttpTask{
			Method:  provider.HttpMethodPost,
			Request: syncTaskEndpoint,
			Token:   env.GetString("PODOPS_API_KEY", ""),
			Payload: &ir,
		}

		err := background().CreateHttpTask(ctx, task)
		if err != nil {
			return err
		}
	}

	// track api access for billing etc
	platform.Meter(ctx, "api.build", "production", p.GUID)

	resp := podops.BuildRequest{
		GUID:         req.GUID,
		FeedURL:      fmt.Sprintf("%s/%s/feed.xml", podops.DefaultStorageEndpoint, req.GUID),
		FeedAliasURL: fmt.Sprintf("%s/s/%s/feed.xml", podops.DefaultEndpoint, p.Name),
	}

	return api.StandardResponse(c, http.StatusCreated, &resp)
}
