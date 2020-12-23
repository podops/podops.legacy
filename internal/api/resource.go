package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/podops/podops/internal/errors"
	"github.com/podops/podops/pkg/metadata"

	"github.com/podops/podops/internal/production"
)

// ResourceEndpoint creates or updates a resource
func ResourceEndpoint(c *gin.Context) {

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
