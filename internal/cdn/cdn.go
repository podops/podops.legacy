package cdn

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/analytics"
	"github.com/podops/podops/internal/api"
	"github.com/podops/podops/pkg/backend"
	"google.golang.org/appengine"
)

// FeedEndpoint handles request for feed.xml by redirecting to the public storage bucket
func FeedEndpoint(c echo.Context) error {

	name := c.Param("name")
	if name == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':name'"))
	}
	prod, err := backend.FindProductionByName(appengine.NewContext(c.Request()), name)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if prod == nil {
		return api.ErrorResponse(c, http.StatusNotFound, fmt.Errorf("can not find '%s/feed.xml'", name))
	}

	redirectTo := fmt.Sprintf("%s/%s/feed.xml", a.StorageEndpoint, prod.GUID)

	// track the event
	analytics.TrackEvent(c.Request(), "cdn", "feed", prod.GUID, 1)

	return c.Redirect(http.StatusTemporaryRedirect, redirectTo)
}
