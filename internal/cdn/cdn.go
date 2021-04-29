package cdn

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/server"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/errordef"
	lp "github.com/podops/podops/internal/platform"
)

// FIXME move this to the caddy handler ?

// FeedEndpoint handles request for feed.xml by redirecting to the public storage bucket
func FeedEndpoint(c echo.Context) error { // FIXME not needed !

	name := c.Param("name")
	if name == "" {
		return server.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}

	prod, err := backend.FindProductionByName(platform.NewHttpContext(c.Request()), name)
	if err != nil {
		return server.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if prod == nil {
		return server.ErrorResponse(c, http.StatusNotFound, errordef.ErrNoSuchProduction)
	}

	redirectTo := fmt.Sprintf("%s/%s/feed.xml", podops.DefaultStorageEndpoint, prod.GUID)

	// track the event
	lp.TrackEvent(c.Request(), "cdn", "feed", prod.GUID, 1)

	return c.Redirect(http.StatusTemporaryRedirect, redirectTo)
}
