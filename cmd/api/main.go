package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/auth"
	"github.com/txsvc/service/pkg/svc"

	"github.com/podops/podops/pkg/api"
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

	// used to secure cookies and sign the JWT token
	secret := env.GetString("MASTER_KEY", "supersecretsecret")

	// create the JWT middleware
	jwt, err := auth.GetSecureJWTMiddleware(env.GetString("REALM", "podops"), secret)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// default public endpoints without authentication
	svc.GET("/", svc.NullEndpoint)

	// Admin Endpoints
	admin := svc.Group(api.AdminNamespacePrefix)
	admin.POST("/token", auth.CreateJWTAuthorizationEndpoint)
	admin.GET("/token", auth.ValidateJWTAuthorizationEndpoint)

	// API endpoints with authentication
	apiEndpoints := svc.SecureGroup(api.NamespacePrefix, jwt.MiddlewareFunc())
	apiEndpoints.POST(api.NewShowRoute, "api.create", api.NewShowEndpoint)

	// add CORS handler, allowing all. See https://github.com/gin-contrib/cors
	svc.Use(cors.Default())

	// add session handler
	store := cookie.NewStore([]byte(secret))
	svc.Use(sessions.Sessions("pops", store))

	// add the service/router to a server on $PORT and launch it. This call BLOCKS !
	svc.Start()
}
