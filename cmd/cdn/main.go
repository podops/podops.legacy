package main

import (
	"log"

	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/pkg/env"
	"github.com/txsvc/platform/v2/provider/google"
	"github.com/txsvc/platform/v2/provider/local"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/podops/podops/internal/cdn"
)

func init() {
	// initialize the platform first
	if !env.Assert("PROJECT_ID") {
		log.Fatal("Missing env variable 'PROJECT_ID'")
	}

	local.InitDefaultProviders()
	p := platform.DefaultPlatform()
	err := p.RegisterProviders(true, google.GoogleErrorReportingConfig)
	if err != nil {
		log.Fatal("error initializing the platform services")
	}

	platform.RegisterPlatform(p) // redundant, but in case we return a copy in the future ...
}

func main() {
	caddycmd.Main()
}
