package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/podops/podops/internal/errors"
	"github.com/podops/podops/internal/production"
	"github.com/podops/podops/pkg/metadata"
)

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

	err := production.CreateResource(appengine.NewContext(c.Request), fmt.Sprintf("%s/%s-%s.yaml", parent, kind, guid), forceFlag, payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	// FIXME validate guid and parent against whatever is in the doc!

	/*
		OLD IMPLEMENTATION
		if resource == "show" {
			var show metadata.Show

			err := c.BindJSON(&show)
			if err != nil {
				HandleError(c, err)
				return
			}

			err = production.CreateResource(appengine.NewContext(c.Request), fmt.Sprintf("%s/show-%s.yaml", parent, guid), forceFlag, &show)
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

			err = production.CreateResource(appengine.NewContext(c.Request), fmt.Sprintf("%s/episode-%s.yaml", parent, guid), forceFlag, &episode)
			if err != nil {
				HandleError(c, err)
				return
			}

		} else {
			HandleError(c, errors.New(fmt.Sprintf("Invalid resource. '%s", resource), http.StatusBadRequest))
			return
		}
	*/

	StandardResponse(c, http.StatusCreated, nil)
}
