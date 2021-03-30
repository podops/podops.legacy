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

	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/auth"
	"github.com/podops/podops/internal/api"
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
	tasks := e.Group(apiv1.TaskNamespacePrefix)
	tasks.POST(apiv1.ImportTask, api.ImportTaskEndpoint)

	// admin endpoints
	admin := e.Group(apiv1.AdminNamespacePrefix)
	admin.POST(apiv1.LoginRequestRoute, auth.LoginRequestEndpoint)
	admin.POST(apiv1.LogoutRequestRoute, auth.LogoutRequestEndpoint)
	admin.GET(apiv1.LoginConfirmationRoute, auth.LoginConfirmationEndpoint)
	admin.POST(apiv1.GetAuthorizationRoute, auth.GetAuthorizationEndpoint)

	// the api endpoints
	apiEndpoints := e.Group(apiv1.NamespacePrefix)
	apiEndpoints.GET(apiv1.ListProductionsRoute, api.ListProductionsEndpoint)
	apiEndpoints.POST(apiv1.ProductionRoute, api.ProductionEndpoint)
	apiEndpoints.GET(apiv1.FindResourceRoute, api.FindResourceEndpoint)
	apiEndpoints.GET(apiv1.GetResourceRoute, api.GetResourceEndpoint)
	apiEndpoints.GET(apiv1.ListResourcesRoute, api.ListResourcesEndpoint)
	apiEndpoints.POST(apiv1.UpdateResourceRoute, api.UpdateResourceEndpoint)
	apiEndpoints.PUT(apiv1.UpdateResourceRoute, api.UpdateResourceEndpoint)
	apiEndpoints.DELETE(apiv1.DeleteResourceRoute, api.DeleteResourceEndpoint)
	apiEndpoints.POST(apiv1.BuildRoute, api.BuildFeedEndpoint)
	apiEndpoints.POST(apiv1.UploadRoute, api.UploadEndpoint)

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
