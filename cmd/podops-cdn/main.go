package main

import (
	"fmt"
	"net/http"

	"github.com/fupas/commons/pkg/env"
	svc "github.com/fupas/platform/pkg/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podops/podops/internal/analytics"
	"github.com/podops/podops/internal/api"
	"github.com/podops/podops/internal/cdn"
)

// ShutdownDelay is the delay before exiting the process
const ShutdownDelay = 10

var (
	// the router instance
	mux                *echo.Echo
	staticFileLocation string
	showPagePath       string
	episodePagePath    string
)

func setup() *echo.Echo {
	// Create a new router instance
	e := echo.New()

	// add and configure the middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// TODO: add/configure e.Use(middleware.Logger())
	// TODO: e.Logger.SetLevel(log.INFO)

	e.GET(api.ShowRoute, RewriteShowHandler)
	e.GET(api.EpisodeRoute, RewriteEpisodeHandler)
	e.GET(api.FeedRoute, cdn.FeedEndpoint)

	// add the routes last
	e.Static("/", staticFileLocation) // serve static files from e.g. ./public

	return e
}

// RewriteShowHandler rewrites requests from /s/:name to /s/_id.html
func RewriteShowHandler(c echo.Context) error {
	if err := c.File(showPagePath); err != nil {
		c.Logger().Error(err)
	}
	// track the event
	analytics.TrackEvent(c.Request(), "podcast", "show", c.Param("name"), 1)

	return nil
}

// RewriteEpisodeHandler rewrites requests from /s/:name/:guid to /e/_id.html
func RewriteEpisodeHandler(c echo.Context) error {
	if err := c.File(episodePagePath); err != nil {
		c.Logger().Error(err)
	}
	// track the event
	analytics.TrackEvent(c.Request(), "podcast", "episode", c.Param("guid"), 1)

	return nil
}

func shutdown(*echo.Echo) {
	// TODO: implement your own stuff here
}

func init() {
	staticFileLocation = env.GetString("STATIC_FILE_LOCATION", "public")
	showPagePath = fmt.Sprintf("%s/s/_id.html", staticFileLocation)
	episodePagePath = fmt.Sprintf("%s/e/_id.html", staticFileLocation)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%s/%d.html", staticFileLocation, code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func main() {
	service := svc.NewServer(setup, shutdown, nil)
	service.StartBlocking()
}

/*


"/a/v1"

"/_a"
"/_t"
"/c"

"/q"

"/s/:name"
"/s/:name/feed.xml"

"/e/:guid"

*/
