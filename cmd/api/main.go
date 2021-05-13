package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/txsvc/platform/v2"
	authapi "github.com/txsvc/platform/v2/pkg/api"
	"github.com/txsvc/platform/v2/pkg/env"
	"github.com/txsvc/platform/v2/pkg/httpserver"
	"github.com/txsvc/platform/v2/provider/google"

	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/graphql"
	"github.com/podops/podops/internal/provider"
)

// ShutdownDelay is the delay before exiting the process
const ShutdownDelay = 10

func setup() *echo.Echo {
	// create a new router instance
	e := echo.New()

	// add and configure the middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig)) // needed for the GraphQL endpoints

	// admin endpoints
	e.GET(apiv1.LoginConfirmationRoute, authapi.LoginConfirmationEndpoint)
	e.POST(apiv1.LogoutRequestRoute, authapi.LogoutRequestEndpoint)

	admin := e.Group(apiv1.AdminNamespacePrefix)
	admin.POST(apiv1.LoginRequestRoute, authapi.LoginRequestEndpoint)
	//admin.POST(apiv1.LoginRequestRoute, hack.HackEndpoint)
	admin.POST(apiv1.GetAuthorizationRoute, authapi.GetAuthorizationEndpoint)

	// FIXME check this !
	//admin.GET(apiv1.LoginConfirmationRoute, authapi.LoginConfirmationEndpoint)
	//admin.POST(apiv1.LogoutRequestRoute, authapi.LogoutRequestEndpoint)

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

	// grapghql endpoints
	gql := e.Group(apiv1.GraphqlNamespacePrefix)
	gql.POST(apiv1.GraphqlRoute, graphql.GraphqlEndpoint())
	gql.GET(apiv1.GraphqlPlaygroundRoute, graphql.GraphqlPlaygroundEndpoint())

	return e
}

func shutdown(*echo.Echo) {
	platform.Close()
}

func init() {
	// initialize the platform first
	if !env.Assert("PROJECT_ID") {
		log.Fatal("Missing env variable 'PROJECT_ID'")
	}
	if !env.Assert("PODOPS_API_KEY") {
		log.Fatal("Missing env variable 'PODOPS_API_KEY'")
	}
	if !env.Assert("LOCATION_ID") {
		log.Fatal("Missing env variable 'LOCATION_ID'")
	}
	if !env.Assert("DEFAULT_QUEUE") {
		log.Fatal("Missing env variable 'DEFAULT_QUEUE'")
	}
	if !env.Assert("EMAIL_DOMAIN") {
		log.Fatal("Missing env variable 'EMAIL_DOMAIN'")
	}
	if !env.Assert("EMAIL_API_KEY") {
		log.Fatal("Missing env variable 'EMAIL_API_KEY'")
	}

	google.InitGoogleCloudPlatformProviders()
	platform.DefaultPlatform().RegisterProviders(true, provider.PodopsAuthConfig)
}

func main() {
	service := httpserver.New(setup, shutdown, nil)
	service.StartBlocking()
}
