package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/env"
	"github.com/txsvc/platform/pkg/http"
	"github.com/txsvc/platform/provider/google"

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
	if !env.Assert("PROJECT_ID") {
		log.Fatal("Missing env variable 'PROJECT_ID'")
	}

	// FIXME InitDefaultPlatform
	p, err := platform.InitPlatform(context.Background(), google.GoogleErrorReportingConfig, platform.DefaultContextConfig)
	if err != nil {
		log.Fatal("error initializing the platform services")
	}
	platform.RegisterPlatform(p)
}

func main() {
	service := http.New(setup, shutdown, nil)
	service.StartBlocking()
}
