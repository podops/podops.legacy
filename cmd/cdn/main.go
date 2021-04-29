package main

import (
	"context"
	"log"

	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	"github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/env"
	"github.com/txsvc/platform/provider/google"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/podops/podops/internal/cdn"
)

func init() {
	// initialize the platform first
	if !env.Assert("PROJECT_ID") {
		log.Fatal("Missing env variable 'PROJECT_ID'")
	}

	er := platform.PlatformOpts{ID: "platform.google.errorreporting", Type: platform.ProviderTypeErrorReporter, Impl: google.NewErrorReporter}
	p, err := platform.InitPlatform(context.Background(), er)
	if err != nil {
		log.Fatal("error initializing the platform services")
	}
	platform.RegisterPlatform(p)
}

func main() {
	caddycmd.Main()
}
