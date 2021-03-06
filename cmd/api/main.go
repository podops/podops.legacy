package main

import (
	"context"
	"log"

	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/platform"
	svc "github.com/fupas/platform/pkg/http"
	gcp "github.com/fupas/platform/provider/google"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podops/podops/internal/api"
	"github.com/podops/podops/pkg/auth"
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
	admin := e.Group(api.AdminNamespacePrefix)
	admin.POST(api.AuthenticationRoute, auth.CreateAuthorizationEndpoint)
	admin.GET(api.AuthenticationRoute, auth.ValidateAuthorizationEndpoint)

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
	// initialize the platform first
	projectID := env.GetString("PROJECT_ID", "")
	if projectID == "" {
		log.Fatal("Missing variable 'PROJECT_ID'")
	}
	serviceName := env.GetString("SERVICE_NAME", "default")

	client, err := platform.NewClient(context.Background(), gcp.NewErrorReporting(context.TODO(), projectID, serviceName))
	if err != nil {
		log.Fatal("error initializing the platform services")
	}
	platform.RegisterGlobally(client)
}

func main() {
	service := svc.NewServer(setup, shutdown, nil)
	service.StartBlocking()
}
