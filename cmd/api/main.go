package main

import (
	svc "github.com/fupas/platform/pkg/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podops/podops/internal/api"
	"github.com/podops/podops/pkg/backend"
)

// ShutdownDelay is the delay before exiting the process
const ShutdownDelay = 10

var (
	// the router instance
	mux *echo.Echo
)

func setup() *echo.Echo {
	// Create a new router instance
	e := echo.New()

	// add and configure the middlewares
	e.Use(middleware.Recover())

	//e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	// TODO: add/configure e.Use(middleware.Logger())
	// TODO: e.Logger.SetLevel(log.INFO)

	// default public endpoints without authentication
	e.GET(api.VersionRoute, api.VersionEndpoint)

	// task endpoints
	tasks := e.Group(api.TaskNamespacePrefix)
	tasks.POST(backend.ImportTask, backend.ImportTaskEndpoint)

	// admin endpoints
	// admin := svc.Group(api.AdminNamespacePrefix)
	//admin.POST(api.AuthenticationRoute, acl.CreateAuthorizationEndpoint)
	//admin.GET(api.AuthenticationRoute, acl.ValidateAuthorizationEndpoint)

	// the api endpoints
	apiEndpoints := e.Group(api.NamespacePrefix)
	apiEndpoints.GET(api.ListProductionsRoute, api.ListProductionsEndpoint)
	apiEndpoints.POST(api.ProductionRoute, api.ProductionEndpoint)
	apiEndpoints.GET(api.GetResourceRoute, api.GetResourceEndpoint)
	apiEndpoints.GET(api.ListResourcesRoute, api.ListResourcesEndpoint)
	apiEndpoints.POST(api.UpdateResourceRoute, api.UpdateResourceEndpoint)
	apiEndpoints.PUT(api.UpdateResourceRoute, api.UpdateResourceEndpoint)
	apiEndpoints.DELETE(api.DeleteResourceRoute, api.DeleteResourceEndpoint)
	apiEndpoints.POST(api.BuildRoute, api.BuildEndpoint)
	apiEndpoints.POST(api.UploadRoute, api.UploadEndpoint)

	return e
}

func shutdown(*echo.Echo) {
	// TODO: implement your own stuff here
}

func init() {
	// TODO implement all the global stuff here
}

func main() {
	service := svc.NewServer(setup, shutdown, nil)
	service.StartBlocking()
}
