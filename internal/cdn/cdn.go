package cdn

/* See the following resourced for reference:

https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers

https://github.com/gin-gonic/gin

https://cloud.google.com/cdn/docs/

*/

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/platform/pkg/platform"
	"github.com/labstack/echo/v4"
	a "github.com/podops/podops/apiv1"
	p "github.com/podops/podops/internal/platform"
	"github.com/podops/podops/pkg/api"
	"github.com/podops/podops/pkg/backend"
	"google.golang.org/appengine"
)

const (
	cacheControl = "public, max-age=1800"
)

var (
	staticFileLocation string
	showPagePath       string
	episodePagePath    string
	bkt                *storage.BucketHandle
)

func init() {
	staticFileLocation = env.GetString("STATIC_FILE_LOCATION", "./public")
	bkt = platform.Storage().Bucket(a.BucketCDN)
}

// RewriteShowHandler rewrites requests from /s/:name to /s/_id.html
func RewriteShowHandler(c echo.Context) error {
	if err := c.File(showPagePath); err != nil {
		c.Logger().Error(err)
	}
	// track the event
	p.TrackEvent(c.Request(), "podcast", "show", c.Param("name"), 1)

	return nil
}

// RewriteEpisodeHandler rewrites requests from /e/:guid to /e/_id.html
func RewriteEpisodeHandler(c echo.Context) error {
	if err := c.File(episodePagePath); err != nil {
		c.Logger().Error(err)
	}
	// track the event
	p.TrackEvent(c.Request(), "podcast", "episode", c.Param("guid"), 1)

	return nil
}

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
	p.TrackEvent(c.Request(), "cdn", "feed", prod.GUID, 1)

	return c.Redirect(http.StatusTemporaryRedirect, redirectTo)
}

// RedirectCDNContentEndpoint serves request for content by redirecting to the public Cloud Storage bucket.
// HEAD, GET are supported operations.
func RedirectCDNContentEndpoint(c echo.Context) error {
	// return an error if the request is anything other than GET/HEAD
	m := c.Request().Method
	if m != "" && m != "GET" && m != "HEAD" {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("received a '%s' request", m))
	}

	guid := c.Param("guid")
	if guid == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected '/:guid/:asset'"))
	}
	asset := c.Param("asset")
	if asset == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected '/:guid/:asset'"))
	}
	rsrc := fmt.Sprintf("%s/%s", guid, asset)

	// handle HEAD request
	if m == "HEAD" {
		// get object attributes, can be cached ...
		obj := bkt.Object(rsrc)
		attr, err := obj.Attrs(appengine.NewContext(c.Request()))

		if err == storage.ErrObjectNotExist {
			return api.ErrorResponse(c, http.StatusNotFound, fmt.Errorf("can not find '%s'", rsrc))
		}

		c.Response().Header().Set("etag", attr.Etag)
		c.Response().Header().Set("accept-ranges", "bytes")
		c.Response().Header().Set("cache-control", cacheControl)
		c.Response().Header().Set("accept-ranges", "bytes")
		c.Response().Header().Set("content-type", attr.ContentType)
		c.Response().Header().Set("content-length", fmt.Sprintf("%d", attr.Size))

		// track the event
		p.TrackEvent(c.Request(), "cdn", "asset", rsrc, 1)

		return c.NoContent(http.StatusOK)
	}

	// track the event
	p.TrackEvent(c.Request(), "cdn", "asset", rsrc, 1)

	// let the storage cdn handle the request
	redirectTo := fmt.Sprintf("%s/%s", a.StorageEndpoint, rsrc)
	return c.Redirect(http.StatusTemporaryRedirect, redirectTo)
}
