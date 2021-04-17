package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/podops/podops/internal/cdn"
)

func main() {
	caddycmd.Main()
}
