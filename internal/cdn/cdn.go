package cdn

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/platform"
)

// FeedEndpoint handles request for feed.xml by redirecting to the public storage bucket
func FeedEndpoint(c echo.Context) error { // FIXME not needed !

	name := c.Param("name")
	if name == "" {
		return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':name'"))
	}

	prod, err := backend.FindProductionByName(platform.NewHttpContext(c), name)
	if err != nil {
		return platform.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if prod == nil {
		return platform.ErrorResponse(c, http.StatusNotFound, fmt.Errorf("can not find '%s/feed.xml'", name))
	}

	redirectTo := fmt.Sprintf("%s/%s/feed.xml", podops.StorageEndpoint, prod.GUID)

	// track the event
	platform.TrackEvent(c.Request(), "cdn", "feed", prod.GUID, 1)

	return c.Redirect(http.StatusTemporaryRedirect, redirectTo)
}
