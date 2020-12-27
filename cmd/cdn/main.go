package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/platform/pkg/platform"

	"github.com/podops/podops/internal/cdn"
)

func init() {
	// setup shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		shutdown()
		os.Exit(1)
	}()
}

func shutdown() {
	platform.Close()
	log.Printf("Exiting ...")
}

func main() {

	// basic http stack config
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(gin.Recovery())

	// the only, catch-all route for https://cdn.podops.dev/*
	r.NoRoute(cdn.ServeContentEndpoint)

	// start the router on port 8080, unless $PORT is set to something else
	r.Run()

}
