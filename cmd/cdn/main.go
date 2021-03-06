package main

import (
	"fmt"
	"net/http"

	"github.com/fupas/commons/pkg/env"
	svc "github.com/fupas/platform/pkg/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podops/podops/internal/api"
	"github.com/podops/podops/internal/cdn"
)

// ShutdownDelay is the delay before exiting the process
const ShutdownDelay = 10

var (
	// the router instance
	mux                *echo.Echo
	staticFileLocation string
)

func setup() *echo.Echo {
	// Create a new router instance
	e := echo.New()

	// hack to get get rid of these 404 in the log
	e.Pre(middleware.Rewrite(map[string]string{
		"/7c054e6693dc/feed.xml": "/s/wizards-magic-sheep/feed.xml",
	}))
	// end hack

	// add and configure the middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(middleware.CSRFWithConfig(middleware.DefaultCSRFConfig))

	// TODO: add/configure e.Use(middleware.Logger())
	// TODO: e.Logger.SetLevel(log.INFO)

	// frontend routes for feed, show & episode
	e.GET(api.ShowRoute, cdn.RewriteShowHandler)
	e.GET(api.EpisodeRoute, cdn.RewriteEpisodeHandler)
	e.GET(api.FeedRoute, cdn.FeedEndpoint)

	// cdn enpoints
	content := e.Group(api.ContentNamespace)
	content.GET(api.DefaultCDNRoute, cdn.RedirectCDNContentEndpoint)
	content.HEAD(api.DefaultCDNRoute, cdn.RedirectCDNContentEndpoint)

	// grapghql
	gql := e.Group(api.GraphqlNamespacePrefix)
	gql.POST(api.GraphqlRoute, api.GetGraphqlEndpoint())
	gql.GET(api.GraphqlPlaygroundRoute, api.GetGraphqlPlaygroundEndpoint())

	// add the routes last
	e.Static("/", staticFileLocation) // serve static files from e.g. ./public

	return e
}

func shutdown(*echo.Echo) {
	// TODO: implement your own stuff here
}

func init() {
	staticFileLocation = env.GetString("STATIC_FILE_LOCATION", "public")
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
