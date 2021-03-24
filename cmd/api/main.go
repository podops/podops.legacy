package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/platform"
	svc "github.com/fupas/platform/pkg/http"
	gcp "github.com/fupas/platform/provider/google"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/api"
	"github.com/podops/podops/pkg/auth"
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

	// task endpoints
	tasks := e.Group(a.TaskNamespacePrefix)
	tasks.POST(a.ImportTask, api.ImportTaskEndpoint)

	// admin endpoints
	admin := e.Group(a.AdminNamespacePrefix)
	admin.POST(a.LoginRequestRoute, auth.LoginRequestEndpoint)
	admin.POST(a.LogoutRequestRoute, auth.LogoutRequestEndpoint)
	admin.GET(a.LoginConfirmationRoute, auth.LoginConfirmationEndpoint)
	admin.POST(a.GetAuthorizationRoute, auth.GetAuthorizationEndpoint)

	// the api endpoints
	apiEndpoints := e.Group(a.NamespacePrefix)
	apiEndpoints.GET(a.ListProductionsRoute, api.ListProductionsEndpoint)
	apiEndpoints.POST(a.ProductionRoute, api.ProductionEndpoint)
	apiEndpoints.GET(a.FindResourceRoute, api.FindResourceEndpoint)
	apiEndpoints.GET(a.GetResourceRoute, api.GetResourceEndpoint)
	apiEndpoints.GET(a.ListResourcesRoute, api.ListResourcesEndpoint)
	apiEndpoints.POST(a.UpdateResourceRoute, api.UpdateResourceEndpoint)
	apiEndpoints.PUT(a.UpdateResourceRoute, api.UpdateResourceEndpoint)
	apiEndpoints.DELETE(a.DeleteResourceRoute, api.DeleteResourceEndpoint)
	apiEndpoints.POST(a.BuildRoute, api.BuildFeedEndpoint)
	apiEndpoints.POST(a.UploadRoute, api.UploadEndpoint)

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
