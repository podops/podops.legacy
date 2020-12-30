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

	"github.com/podops/podops/internal/api"
	"github.com/podops/podops/internal/resources"
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
	admin.POST(api.AuthenticationRoute, auth.CreateJWTAuthorizationEndpoint)
	admin.GET(api.AuthenticationRoute, auth.ValidateJWTAuthorizationEndpoint)

	// Task Endpoints
	tasks := svc.Group(api.TaskNamespacePrefix)
	tasks.POST(resources.ImportTask, resources.ImportTaskEndpoint)

	// API endpoints with authentication
	apiEndpoints := svc.SecureGroup(api.NamespacePrefix, jwt.MiddlewareFunc())
	apiEndpoints.GET(api.ListRoute, "api.view", api.ListProductionsEndpoint)
	apiEndpoints.POST(api.ProductionRoute, "api.create", api.ProductionEndpoint)
	apiEndpoints.POST(api.ResourceRoute, "api.create,api.update", api.ResourceEndpoint) // creates a resource, fails if it already exists
	apiEndpoints.PUT(api.ResourceRoute, "api.update", api.ResourceEndpoint)             // updates a resource, fails if it does NOT exist
	apiEndpoints.POST(api.BuildRoute, "api.update", api.BuildEndpoint)
	apiEndpoints.POST(api.UploadRoute, "api.update", api.UploadEndpoint)

	// add CORS handler, allowing all. See https://github.com/gin-contrib/cors
	svc.Use(cors.Default())

	// add session handler
	store := cookie.NewStore([]byte(secret))
	svc.Use(sessions.Sessions("pops", store))

	// add the service/router to a server on $PORT and launch it. This call BLOCKS !
	svc.Start()
}
