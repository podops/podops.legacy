package api

import (
	"fmt"
	"net/http"

	"github.com/fupas/commons/pkg/util"
	"github.com/labstack/echo/v4"
	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/analytics"
	"github.com/podops/podops/pkg/auth"
	"github.com/podops/podops/pkg/backend"
	"google.golang.org/appengine"
)

// GetResourceEndpoint returns a resource
func GetResourceEndpoint(c echo.Context) error {
	if status, err := auth.Authorized(c, "ROLES"); err != nil {
		return ErrorResponse(c, status, err)
	}

	prod := c.Param("prod")
	if prod == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':prod"))
	}
	kind := c.Param("kind")
	if kind == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':kind"))
	}
	guid := c.Param("id")
	if guid == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':id"))
	}

	// FIXME prod, kind are ignored, assumption is that guid is globally unique ...

	resource, err := backend.GetResourceContent(appengine.NewContext(c.Request()), guid)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	analytics.TrackEvent(c.Request(), "api", "rsrc_get", fmt.Sprintf("%s/%s/%s", prod, kind, guid), 1)

	if resource == nil {
		return StandardResponse(c, http.StatusNotFound, nil)
	}
	return StandardResponse(c, http.StatusOK, resource)

}

// ListResourcesEndpoint returns a list of resources
func ListResourcesEndpoint(c echo.Context) error {
	if status, err := auth.Authorized(c, "ROLES"); err != nil {
		return ErrorResponse(c, status, err)
	}

	prod := c.Param("prod")
	if prod == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':prod"))
	}
	kind := c.Param("kind")
	if kind == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':kind"))
	}

	l, err := backend.ListResources(appengine.NewContext(c.Request()), prod, kind)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	analytics.TrackEvent(c.Request(), "api", "rsrc_list", fmt.Sprintf("%s/%s", prod, kind), 1)

	return StandardResponse(c, http.StatusOK, &a.ResourceList{Resources: l})
}

// UpdateResourceEndpoint creates or updates a resource
func UpdateResourceEndpoint(c echo.Context) error {
	if status, err := auth.Authorized(c, "ROLES"); err != nil {
		return ErrorResponse(c, status, err)
	}

	prod := c.Param("prod")
	if prod == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':prod"))
	}
	kind := c.Param("kind")
	if kind == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':kind"))
	}
	guid := c.Param("id")
	if guid == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':id"))
	}

	forceFlag := false
	if c.QueryParam("f") == "true" {
		forceFlag = true
	}

	var payload interface{}
	ctx := appengine.NewContext(c.Request())
	location := fmt.Sprintf("%s/%s-%s.yaml", prod, kind, guid)

	if kind == a.ResourceShow {
		var show *a.Show = new(a.Show)

		if err := c.Bind(show); err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, err)
		}
		payload = &show

		if prod != show.GUID() {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(":prod and GUID do not match. expected '%s', got '%s'", prod, show.GUID()))
		}

		// update the PRODUCTION entry based on resource
		p, err := backend.GetProduction(ctx, show.GUID())
		if err != nil {
			return ErrorResponse(c, http.StatusNotFound, err)
		}

		// the attributes we copy from the .yaml
		p.Title = show.Description.Title
		p.Summary = show.Description.Summary
		p.Updated = util.Timestamp()

		if err := backend.UpdateProduction(ctx, p); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.EnsureAsset(ctx, show.GUID(), &show.Image); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.UpdateShow(ctx, location, show); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, err)
		}

	} else if kind == a.ResourceEpisode {
		var episode *a.Episode = new(a.Episode)

		if err := c.Bind(episode); err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, err)
		}
		payload = &episode

		if prod != episode.ParentGUID() {
			return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(":prod and GUID do not match. expected '%s', got '%s'", prod, episode.ParentGUID()))
		}

		// ensure images and media files
		if err := backend.EnsureAsset(ctx, episode.ParentGUID(), &episode.Image); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.EnsureAsset(ctx, episode.ParentGUID(), &episode.Enclosure); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.UpdateEpisode(ctx, location, episode); err != nil {
			return ErrorResponse(c, http.StatusBadRequest, err)
		}
	} else {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("unsupported kind '%s'", kind))
	}

	createFlag := true // POST
	action := "rsrc_create"

	if c.Request().Method == "PUT" {
		createFlag = false
		action = "rsrc_update"
	}
	if err := backend.WriteResourceContent(ctx, location, createFlag, forceFlag, payload); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	analytics.TrackEvent(c.Request(), "api", action, fmt.Sprintf("%s/%s/%s", prod, kind, guid), 1)

	return StandardResponse(c, http.StatusCreated, nil)
}

// DeleteResourceEndpoint deletes a resource and its .yaml file
func DeleteResourceEndpoint(c echo.Context) error {
	if status, err := auth.Authorized(c, "ROLES"); err != nil {
		return ErrorResponse(c, status, err)
	}

	prod := c.Param("prod")
	if prod == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':prod"))
	}
	kind := c.Param("kind")
	if kind == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':kind"))
	}
	guid := c.Param("id")
	if guid == "" {
		return ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':id"))
	}

	// FIXME prod, kind are ignored, assumption is that guid is globally unique ...

	if err := backend.DeleteResource(appengine.NewContext(c.Request()), guid); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	analytics.TrackEvent(c.Request(), "api", "rsrc_delete", fmt.Sprintf("%s/%s/%s", prod, kind, guid), 1)

	return c.NoContent(http.StatusNoContent)
}
