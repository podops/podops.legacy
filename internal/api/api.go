package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/podops/podops/internal/errors"
	"github.com/podops/podops/internal/resources"
	t "github.com/podops/podops/internal/types"
	"github.com/podops/podops/pkg/metadata"
)

const (
	// AdminNamespacePrefix namespace for internal admin endpoints
	AdminNamespacePrefix = "/_a"
	// NamespacePrefix namespace for the CLI. Should not be used directly.
	NamespacePrefix = "/a/v1"

	// All the API & CLI endpoint routes

	// AuthenticationRoute is used to create and verify a token
	AuthenticationRoute = "/token"
	// ProductionRoute route to ProductionEndpoint
	ProductionRoute = "/new"
	// ResourceRoute route to ResourceEndpoint
	ResourceRoute = "/update/:parent/:kind/:id"
)

// ProductionEndpoint creates an new show and does all the background setup
func ProductionEndpoint(c *gin.Context) {
	var req t.ProductionRequest

	err := c.BindJSON(&req)
	if err != nil {
		HandleError(c, err)
		return
	}

	// create a show
	// FIXME: verify && cleanup the name. Should follow Domain name conventions.
	showName := strings.ToLower(strings.TrimSpace(req.Name))
	p, err := resources.CreateProduction(appengine.NewContext(c.Request), showName, req.Title, req.Summary)
	if err != nil {
		HandleError(c, err)
		return
	}

	// send the GUID and canonical name back
	resp := t.ProductionResponse{
		Name: p.Name,
		GUID: p.GUID,
	}
	StandardResponse(c, http.StatusCreated, &resp)
}

// ResourceEndpoint creates or updates a resource
func ResourceEndpoint(c *gin.Context) {

	parent := c.Param("parent")
	if parent == "" {
		HandleError(c, errors.New("Invalid route. Expected ':parent", http.StatusBadRequest))
		return
	}
	kind := c.Param("kind")
	if kind == "" {
		HandleError(c, errors.New("Invalid route. Expected ':kind", http.StatusBadRequest))
		return
	}
	guid := c.Param("id")
	if guid == "" {
		HandleError(c, errors.New("Invalid route. Expected ':id", http.StatusBadRequest))
		return
	}

	//force := c.DefaultQuery("force", "false")
	forceFlag := true
	var payload interface{}

	if kind == "show" {
		var show metadata.Show

		err := c.BindJSON(&show)
		if err != nil {
			HandleError(c, err)
			return
		}
		payload = &show
	} else if kind == "episode" {
		var episode metadata.Episode

		err := c.BindJSON(&episode)
		if err != nil {
			HandleError(c, err)
			return
		}
		payload = &episode
	} else {
		HandleError(c, errors.New(fmt.Sprintf("Invalid resource. '%s", kind), http.StatusBadRequest))
		return
	}

	err := resources.CreateResource(appengine.NewContext(c.Request), fmt.Sprintf("%s/%s-%s.yaml", parent, kind, guid), forceFlag, payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	StandardResponse(c, http.StatusCreated, nil)
}
