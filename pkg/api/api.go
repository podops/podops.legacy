package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/platform"
)

const (
	// NamespacePrefix namespace for the client and CLI
	NamespacePrefix = "/a/v1"
	// GraphqlNamespacePrefix namespace for the GraphQL endpoints
	GraphqlNamespacePrefix = "/q"
	// AdminNamespacePrefix namespace for internal admin endpoints
	AdminNamespacePrefix = "/_a"
	// TaskNamespacePrefix namespace for internal Cloud Task callbacks
	TaskNamespacePrefix = "/_t"
	// ContentNamespace namespace for thr CDN
	ContentNamespace = "/c"

	// All the API & CLI endpoint routes

	// VersionRoute route to VersionEndpoint
	VersionRoute = "/version"

	// AuthenticationRoute is used to create and verify a token
	AuthenticationRoute = "/token"

	// ProductionRoute route to ProductionEndpoint
	ProductionRoute = "/production"

	// ListProductionsRoute route to ListProductionsEndpoint
	ListProductionsRoute = "/productions"

	// GetResourceRoute route to ResourceEndpoint
	GetResourceRoute = "/resource/:prod/:kind/:id"

	// ListResourcesRoute route to ResourceEndpoint GET
	ListResourcesRoute = "/resource/:prod/:kind"

	// UpdateResourceRoute route to ResourceEndpoint POST,PUT
	UpdateResourceRoute = "/resource/:prod/:kind/:id"

	// DeleteResourceRoute route to ResourceEndpoint
	DeleteResourceRoute = "/resource/:prod/:kind/:id"

	// BuildRoute route to BuildEndpoint
	BuildRoute = "/build"

	// UploadRoute route to UploadEndpoint
	UploadRoute = "/upload/:prod"

	// ShowRoute route to show.json
	ShowRoute = "/s/:name"

	// EpisodeRoute route to show.json
	EpisodeRoute = "/e/:guid"

	// FeedRoute route to feed.xml
	FeedRoute = "/s/:name/feed.xml"

	// DefaultCDNRoute route to /:guid/:asset
	DefaultCDNRoute = "/:guid/:asset"

	// GraphqlRoute route to GraphqlEndpoint
	GraphqlRoute = "/query"

	// GraphqlPlaygroundRoute route to GraphqlPlaygroundEndpoint
	GraphqlPlaygroundRoute = "/playground"
)

// StandardResponse is the default way to respond to API requests
func StandardResponse(c echo.Context, status int, res interface{}) error {
	if res == nil {
		resp := a.StatusObject{
			Status:  status,
			Message: fmt.Sprintf("status: %d", status),
		}
		return c.JSON(status, &resp)
	} else {
		return c.JSON(status, res)
	}
}

// ErrorResponse reports the error and responds with an ErrorObject
func ErrorResponse(c echo.Context, status int, err error) error {
	var resp a.StatusObject

	// send the error to Google Error Reporting
	platform.ReportError(err)

	if err == nil {
		resp = a.NewStatus(http.StatusInternalServerError, fmt.Sprintf("status: %d", status))
	} else {
		resp = a.NewErrorStatus(status, err)
	}
	return c.JSON(status, &resp)
}

// VersionEndpoint returns the current API version
func VersionEndpoint(c echo.Context) error {
	return c.JSON(http.StatusOK, gin.H{"version": a.VersionString, "major": a.MajorVersion, "minor": a.MinorVersion, "fix": a.FixVersion, "namespace": a.Version})
}
