package cdn

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/api"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/errordef"
)

// FIXME move this to the caddy handler ?

// FeedEndpoint handles request for feed.xml by redirecting to the public storage bucket
func FeedEndpoint(c echo.Context) error { // FIXME not needed !

	name := c.Param("name")
	if name == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}

	prod, err := backend.FindProductionByName(platform.NewHttpContext(c.Request()), name)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if prod == nil {
		return api.ErrorResponse(c, http.StatusNotFound, errordef.ErrNoSuchProduction)
	}

	redirectTo := fmt.Sprintf("%s/%s/feed.xml", podops.DefaultStorageEndpoint, prod.GUID)

	// track api access for billing etc
	platform.Logger("metrics").Log("cdn.feed", "production", prod.GUID)

	return c.Redirect(http.StatusTemporaryRedirect, redirectTo)
}
