package api

import (
	"fmt"
	"net/http"

	"github.com/fupas/commons/pkg/util"
	"github.com/labstack/echo/v4"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/platform"
	"github.com/podops/podops/pkg/backend"
)

// FindResourceEndpoint returns a resource
func FindResourceEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c)

	guid := c.Param("id")
	if guid == "" {
		return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route"))
	}

	if err := AuthorizeAccessResource(ctx, c, scopeResourceRead, guid); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	resource, err := backend.GetResourceContent(ctx, guid)
	if err != nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// FIXME verify that we actually are the owner !!!

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", "rsrc_find", guid, 1)

	if resource == nil {
		return platform.StandardResponse(c, http.StatusNotFound, nil)
	}
	return platform.StandardResponse(c, http.StatusOK, resource)

}

// GetResourceEndpoint returns a resource
func GetResourceEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c)

	prod := c.Param("prod")
	kind := c.Param("kind")
	guid := c.Param("id")

	if !validateNotEmpty(prod, kind, guid) {
		return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route"))
	}

	if err := AuthorizeAccessResource(ctx, c, scopeResourceRead, guid); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	// FIXME prod, kind are ignored, assumption is that guid is globally unique ...
	resource, err := backend.GetResourceContent(ctx, guid)
	if err != nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", "rsrc_get", fmt.Sprintf("%s/%s/%s", prod, kind, guid), 1)

	if resource == nil {
		return platform.StandardResponse(c, http.StatusNotFound, nil)
	}
	return platform.StandardResponse(c, http.StatusOK, resource)

}

// ListResourcesEndpoint returns a list of resources
func ListResourcesEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c)

	prod := c.Param("prod")
	kind := c.Param("kind")

	if !validateNotEmpty(prod, kind) {
		return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route"))
	}

	if err := AuthorizeAccessProduction(ctx, c, scopeResourceRead, prod); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	l, err := backend.ListResources(ctx, prod, kind)
	if err != nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", "rsrc_list", fmt.Sprintf("%s/%s", prod, kind), 1)

	return platform.StandardResponse(c, http.StatusOK, &podops.ResourceList{Resources: l})
}

// UpdateResourceEndpoint creates or updates a resource
func UpdateResourceEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c)

	prod := c.Param("prod")
	kind := c.Param("kind")
	guid := c.Param("id")

	if !validateNotEmpty(prod, kind, guid) {
		return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route"))
	}

	if err := AuthorizeAccessResource(ctx, c, scopeResourceWrite, guid); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	forceFlag := false
	if c.QueryParam("f") == "true" {
		forceFlag = true
	}

	var payload interface{}
	location := fmt.Sprintf("%s/%s-%s.yaml", prod, kind, guid)

	if kind == podops.ResourceShow {
		var show *podops.Show = new(podops.Show)

		if err := c.Bind(show); err != nil {
			return platform.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		payload = &show

		if prod != show.GUID() {
			return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(":prod and GUID do not match. expected '%s', got '%s'", prod, show.GUID()))
		}

		// update the PRODUCTION entry based on resource
		p, err := backend.GetProduction(ctx, show.GUID())
		if err != nil {
			return platform.ErrorResponse(c, http.StatusNotFound, err)
		}

		// the attributes we copy from the .yaml
		p.Title = show.Description.Title
		p.Summary = show.Description.Summary
		p.Updated = util.Timestamp()

		if err := backend.UpdateProduction(ctx, p); err != nil {
			return platform.ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.EnsureAsset(ctx, show.GUID(), &show.Image); err != nil {
			return platform.ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.UpdateShow(ctx, location, show); err != nil {
			return platform.ErrorResponse(c, http.StatusBadRequest, err)
		}

	} else if kind == podops.ResourceEpisode {
		var episode *podops.Episode = new(podops.Episode)

		if err := c.Bind(episode); err != nil {
			return platform.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		payload = &episode

		if prod != episode.Parent() {
			return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(":prod and GUID do not match. expected '%s', got '%s'", prod, episode.Parent()))
		}

		// ensure images and media files
		if err := backend.EnsureAsset(ctx, episode.Parent(), &episode.Image); err != nil {
			return platform.ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.EnsureAsset(ctx, episode.Parent(), &episode.Enclosure); err != nil {
			return platform.ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.UpdateEpisode(ctx, location, episode); err != nil {
			return platform.ErrorResponse(c, http.StatusBadRequest, err)
		}
	} else {
		return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("unsupported kind '%s'", kind))
	}

	createFlag := true // POST
	action := "rsrc_create"

	if c.Request().Method == "PUT" {
		createFlag = false
		action = "rsrc_update"
	}
	if err := backend.WriteResourceContent(ctx, location, createFlag, forceFlag, payload); err != nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", action, fmt.Sprintf("%s/%s/%s", prod, kind, guid), 1)

	return platform.StandardResponse(c, http.StatusCreated, nil)
}

// DeleteResourceEndpoint deletes a resource and its .yaml file
func DeleteResourceEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c)

	prod := c.Param("prod")
	kind := c.Param("kind")
	guid := c.Param("id")

	if !validateNotEmpty(prod, kind, guid) {
		return platform.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route"))
	}

	if err := AuthorizeAccessResource(ctx, c, scopeResourceWrite, guid); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	// FIXME prod, kind are ignored, assumption is that guid is globally unique ...
	if err := backend.DeleteResource(ctx, guid); err != nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.TrackEvent(c.Request(), "api", "rsrc_delete", fmt.Sprintf("%s/%s/%s", prod, kind, guid), 1)

	return c.NoContent(http.StatusNoContent)
}
