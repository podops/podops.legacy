package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/podops/podops/internal/config"
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
	// ListRoute route to ListProductionsEndpoint
	ListRoute = "/list"
	// BuildRoute route to BuildEndpoint
	BuildRoute = "/build"
)

// ProductionEndpoint creates an new show and does all the background setup
func ProductionEndpoint(c *gin.Context) {
	var req t.ProductionRequest

	err := c.BindJSON(&req)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	clientID, err := GetClientID(c)
	if err != nil || clientID == "" {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	// create a show
	// FIXME: verify && cleanup the name. Should follow Domain name conventions.
	showName := strings.ToLower(strings.TrimSpace(req.Name))
	p, err := resources.CreateProduction(appengine.NewContext(c.Request), showName, req.Title, req.Summary, clientID)
	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	// send the GUID and canonical name back
	resp := t.ProductionResponse{
		Name: p.Name,
		GUID: p.GUID,
	}
	StandardResponse(c, http.StatusCreated, &resp)
}

// ListProductionsEndpoint creates an new show and does all the background setup
func ListProductionsEndpoint(c *gin.Context) {

	clientID, err := GetClientID(c)
	if err != nil || clientID == "" {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	productions, err := resources.FindProductionsByOwner(appengine.NewContext(c.Request), clientID)
	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	resp := t.ProductionsResponse{}
	resp.List = make([]t.ProductionDetails, len(productions))
	for i, p := range productions {
		resp.List[i].GUID = p.GUID
		resp.List[i].Name = p.Name
		resp.List[i].Title = p.Title
	}

	StandardResponse(c, http.StatusOK, &resp)
}

// ResourceEndpoint creates or updates a resource
func ResourceEndpoint(c *gin.Context) {

	parent := c.Param("parent")
	if parent == "" {
		HandleError(c, http.StatusBadRequest, fmt.Errorf("resource: invalid route, expected ':parent"))
		return
	}
	kind := c.Param("kind")
	if kind == "" {
		HandleError(c, http.StatusBadRequest, fmt.Errorf("resource: invalid route, expected ':kind"))
		return
	}
	guid := c.Param("id")
	if guid == "" {
		HandleError(c, http.StatusBadRequest, fmt.Errorf("resource: invalid route, expected ':id"))
		return
	}

	forceFlag := false
	if c.DefaultQuery("f", "false") == "true" {
		forceFlag = true
	}

	var payload interface{}

	if kind == "show" {
		var show metadata.Show

		err := c.BindJSON(&show)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}
		payload = &show
	} else if kind == "episode" {
		var episode metadata.Episode

		err := c.BindJSON(&episode)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, err)
			return
		}
		payload = &episode
	} else {
		HandleError(c, http.StatusBadRequest, fmt.Errorf("resource: invalid kind '%s", kind))
		return
	}

	createFlag := true // POST
	if c.Request.Method == "PUT" {
		createFlag = false
	}
	err := resources.WriteResource(appengine.NewContext(c.Request), fmt.Sprintf("%s/%s-%s.yaml", parent, kind, guid), createFlag, forceFlag, payload)
	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	StandardResponse(c, http.StatusCreated, nil)
}

// BuildEndpoint starts the build of the feed
func BuildEndpoint(c *gin.Context) {
	var req t.BuildRequest

	err := c.BindJSON(&req)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	clientID, err := GetClientID(c)
	if err != nil || clientID == "" {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	ctx := appengine.NewContext(c.Request)

	p, err := resources.GetProduction(ctx, req.GUID)
	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}
	if p == nil {
		HandleError(c, http.StatusBadRequest, fmt.Errorf("resource show: invalid guid '%s'", req.GUID))
		return
	}
	if p.Owner != clientID {
		// FIXME this is a simplification in the abscence of proper ACLs
		HandleError(c, http.StatusBadRequest, fmt.Errorf("build: user '%s' not allowd to access production '%s'", clientID, req.GUID))
		return
	}

	// start the build
	err = resources.Build(ctx, req.GUID, false) // FIXME make this async, make validateOnly a flag
	if err != nil {
		HandleError(c, http.StatusBadRequest, fmt.Errorf("build: error building production '%s'", req.GUID))
		return
	}

	resp := t.BuildResponse{
		GUID: req.GUID,
		URL:  fmt.Sprintf("%s/%s/feed.xml", config.DefaultCDNEndpoint, req.GUID),
	}

	StandardResponse(c, http.StatusCreated, &resp)
}
