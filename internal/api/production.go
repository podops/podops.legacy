package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/podops/podops/internal/errors"
	"github.com/podops/podops/pkg/metadata"

	"github.com/podops/podops/internal/production"
)

type (
	// ProductionRequest defines the request
	ProductionRequest struct {
		Name    string `json:"name" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Summary string `json:"summary" binding:"required"`
	}

	// ProductionResponse defines the request
	ProductionResponse struct {
		Name string `json:"name" binding:"required"`
		GUID string `json:"guid" binding:"required"`
	}
)

// CreateProductionEndpoint creates an new show and does all the background setup
func CreateProductionEndpoint(c *gin.Context) {
	var req ProductionRequest

	err := c.BindJSON(&req)
	if err != nil {
		HandleError(c, err)
		return
	}

	// create a show
	showName := strings.ToLower(strings.TrimSpace(req.Name)) // FIXME: verify && cleanup the name. Should follow Domain name conventions.

	p, err := production.CreateProduction(appengine.NewContext(c.Request), showName, req.Title, req.Summary)
	if err != nil {
		HandleError(c, err)
		return
	}

	// just send the GUID and canonical name back
	resp := ProductionResponse{
		Name: p.Name,
		GUID: p.GUID,
	}
	StandardResponse(c, http.StatusCreated, &resp)
}

// CreateEndpoint creates a resource
func CreateEndpoint(c *gin.Context) {

	guid := c.Param("id")
	if guid == "" {
		HandleError(c, errors.New("Invalid route. Expected ':id", http.StatusBadRequest))
		return
	}
	resource := c.Param("rsrc")
	if resource == "" {
		HandleError(c, errors.New("Invalid route. Expected ':rsrc", http.StatusBadRequest))
		return
	}
	//force := c.DefaultQuery("force", "false")
	forceFlag := true

	if resource == "show" {
		var show metadata.Show

		err := c.BindJSON(&show)
		if err != nil {
			HandleError(c, err)
			return
		}

		err = production.CreateResource(appengine.NewContext(c.Request), fmt.Sprintf("%s/show-%s.yaml", guid, guid), forceFlag, &show)
		if err != nil {
			HandleError(c, err)
			return
		}

	} else if resource == "episode" {
		var episode metadata.Episode

		err := c.BindJSON(&episode)
		if err != nil {
			HandleError(c, err)
			return
		}

		err = production.CreateResource(appengine.NewContext(c.Request), fmt.Sprintf("%s/episode-%s.yaml", guid, episode.Metadata.Labels[metadata.LabelGUID]), forceFlag, &episode)
		if err != nil {
			HandleError(c, err)
			return
		}

	} else {
		HandleError(c, errors.New(fmt.Sprintf("Invalid resource. '%s", resource), http.StatusBadRequest))
		return
	}

	StandardResponse(c, http.StatusCreated, nil)
}

// UpdateEndpoint creates a resource
func UpdateEndpoint(c *gin.Context) {
	fmt.Println("/update")
	StandardResponse(c, http.StatusAccepted, nil)
}
