package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/fupas/platform"
	gcp "github.com/fupas/platform/provider/google"
	"github.com/txsvc/spa/pkg/env"
	"github.com/txsvc/spa/pkg/server"

	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/cdn"
)

// ShutdownDelay is the delay before exiting the process
const ShutdownDelay = 10

func setup() *echo.Echo {
	// create a new router instance
	e := echo.New()

	// add and configure the middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	// check for being alive
	e.GET(apiv1.CheckAliveRoute, apiv1.CheckAliveEndpoint)

	// webhook endpoints
	webhook := e.Group(apiv1.WebhookNamespacePrefix)
	webhook.POST(apiv1.ImportTask, cdn.ImportTaskEndpoint)
	webhook.POST(apiv1.SyncTask, cdn.SyncTaskEndpoint)
	webhook.DELETE(apiv1.DeleteTask, cdn.DeleteTaskEndpoint)
	webhook.POST(apiv1.UploadRoute, cdn.UploadEndpoint)

	// redirect to the real feed.xml path
	e.GET(apiv1.FeedRoute, cdn.FeedEndpoint)

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
	serviceName := env.GetString("SERVICE_NAME", "api")

	client, err := platform.NewClient(context.Background(), gcp.NewErrorReporting(context.TODO(), projectID, serviceName))
	if err != nil {
		log.Fatal("error initializing the platform services")
	}
	platform.RegisterGlobally(client)
}

func main() {
	service := server.New(setup, shutdown, nil)
	service.StartBlocking()
}
