package main

import (
	"fmt"
	"net/http"

	"github.com/fupas/commons/pkg/env"
	svc "github.com/fupas/platform/pkg/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podops/podops/internal/cdn"
)

// ShutdownDelay is the delay before exiting the process
const ShutdownDelay = 10

// the router instance
var mux *echo.Echo
var staticFileLocation string = env.GetString("STATIC_FILE_LOCATION", "public")

func setup() *echo.Echo {
	// Create a new router instance
	e := echo.New()

	// add and configure the middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// TODO: add/configure e.Use(middleware.Logger())
	// TODO: e.Logger.SetLevel(log.INFO)

	e.GET("/s/:name", rewriteShow)

	//e.GET("/s/:name/:guid", episode)
	e.GET("/s/:name/feed.xml", cdn.FeedEndpoint)

	// add the routes last
	e.Static("/", staticFileLocation) // serve static files from e.g. ./public

	return e
}

func rewriteShow(c echo.Context) error {
	//name := c.Param("name")
	page := fmt.Sprintf("%s/s/_id.html", staticFileLocation)
	if err := c.File(page); err != nil {
		c.Logger().Error(err)
	}
	return nil
}

func episode(c echo.Context) error {
	name := c.Param("name")
	guid := c.Param("guid")

	return c.String(http.StatusOK, name+"-"+guid)
}

func shutdown(*echo.Echo) {
	// TODO: implement your own stuff here
}

func init() {
	// TODO: initialize everything global here
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
