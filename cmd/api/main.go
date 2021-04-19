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
	"github.com/podops/podops/graphql"
)

// ShutdownDelay is the delay before exiting the process
const ShutdownDelay = 10

func setup() *echo.Echo {
	// create a new router instance
	e := echo.New()

	// add and configure the middlewares
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(middleware.Recover())

	// admin endpoints
	admin := e.Group(apiv1.AdminNamespacePrefix)
	admin.POST(apiv1.LoginRequestRoute, auth.LoginRequestEndpoint)
	admin.POST(apiv1.LogoutRequestRoute, auth.LogoutRequestEndpoint)
	admin.GET(apiv1.LoginConfirmationRoute, auth.LoginConfirmationEndpoint)
	admin.POST(apiv1.GetAuthorizationRoute, auth.GetAuthorizationEndpoint)

	// api endpoints
	apiEndpoints := e.Group(apiv1.NamespacePrefix)
	apiEndpoints.GET(apiv1.ListProductionsRoute, apiv1.ListProductionsEndpoint)
	apiEndpoints.POST(apiv1.ProductionRoute, apiv1.ProductionEndpoint)
	apiEndpoints.GET(apiv1.FindResourceRoute, apiv1.FindResourceEndpoint)
	apiEndpoints.GET(apiv1.GetResourceRoute, apiv1.GetResourceEndpoint)
	apiEndpoints.GET(apiv1.ListResourcesRoute, apiv1.ListResourcesEndpoint)
	apiEndpoints.POST(apiv1.UpdateResourceRoute, apiv1.UpdateResourceEndpoint)
	apiEndpoints.PUT(apiv1.UpdateResourceRoute, apiv1.UpdateResourceEndpoint)
	apiEndpoints.DELETE(apiv1.DeleteResourceRoute, apiv1.DeleteResourceEndpoint)
	apiEndpoints.POST(apiv1.BuildRoute, apiv1.BuildFeedEndpoint)
	apiEndpoints.POST(apiv1.UploadRoute, apiv1.UploadEndpoint)

	// grapghql endpoints
	gql := e.Group(apiv1.GraphqlNamespacePrefix)
	gql.POST(apiv1.GraphqlRoute, graphql.GraphqlEndpoint())
	gql.GET(apiv1.GraphqlPlaygroundRoute, graphql.GraphqlPlaygroundEndpoint())

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
