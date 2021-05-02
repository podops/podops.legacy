package apiv1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/pkg/api"
	"github.com/txsvc/platform/v2/pkg/timestamp"
	"github.com/txsvc/platform/v2/pkg/validate"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/messagedef"
)

// FindResourceEndpoint returns a resource
func FindResourceEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c.Request())

	guid := c.Param("id")
	if guid == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}

	if err := AuthorizeAccessResource(ctx, c, ScopeResourceRead, guid); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	resource, err := backend.GetResourceContent(ctx, guid)
	if err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}
	if resource == nil {
		return api.StandardResponse(c, http.StatusNotFound, nil)
	}

	// track api access for billing etc
	platform.Meter(ctx, "api.resource.find", "resource", guid)

	return api.StandardResponse(c, http.StatusOK, resource)
}

// GetResourceEndpoint returns a resource
func GetResourceEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c.Request())

	prod := c.Param("prod")
	kind := c.Param("kind")
	guid := c.Param("id")

	if !validate.NotEmpty(prod, kind, guid) {
		return api.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}

	if err := AuthorizeAccessResource(ctx, c, ScopeResourceRead, guid); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	// FIXME prod, kind are ignored, assumption is that guid is globally unique ...
	resource, err := backend.GetResourceContent(ctx, guid)
	if err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	if resource == nil {
		return api.StandardResponse(c, http.StatusNotFound, nil)
	}

	// track api access for billing etc
	platform.Meter(ctx, "api.resource.get", "production", prod, "resource", guid, "kind", kind)

	return api.StandardResponse(c, http.StatusOK, resource)
}

// ListResourcesEndpoint returns a list of resources
func ListResourcesEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c.Request())

	prod := c.Param("prod")
	kind := c.Param("kind")

	if !validate.NotEmpty(prod, kind) {
		return api.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}

	if err := AuthorizeAccessProduction(ctx, c, ScopeResourceRead, prod); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	l, err := backend.ListResources(ctx, prod, kind)
	if err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.Meter(ctx, "api.resource.list", "production", prod, "kind", kind)

	return api.StandardResponse(c, http.StatusOK, &podops.ResourceList{Resources: l})
}

// UpdateResourceEndpoint creates or updates a resource
func UpdateResourceEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c.Request())
	createFlag := true // c.Request().Method == POST, default
	action := "api.resource.create"

	prod := c.Param("prod")
	kind := c.Param("kind")
	guid := c.Param("id")

	forceFlag := false
	if strings.ToLower(c.QueryParam("f")) == "true" {
		forceFlag = true
	}

	if c.Request().Method == "PUT" {
		createFlag = false
		action = "api.resource.update"
	}

	if !validate.NotEmpty(prod, kind, guid) {
		return api.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}

	if createFlag {
		// this assumes that the resource does not exist i.e. we only validate access to the production
		if err := AuthorizeAccessProduction(ctx, c, ScopeResourceWrite, prod); err != nil {
			return api.ErrorResponse(c, http.StatusUnauthorized, err)
		}
	} else {
		// we assume the resource already exists and we can validate guid and prod
		if err := AuthorizeAccessResource(ctx, c, ScopeResourceWrite, guid); err != nil {
			return api.ErrorResponse(c, http.StatusUnauthorized, err)
		}
	}

	var payload interface{}
	location := fmt.Sprintf("%s/%s-%s.yaml", prod, kind, guid)

	if kind == podops.ResourceShow {
		var show *podops.Show = new(podops.Show) // FIXME change this !

		if err := c.Bind(show); err != nil {
			return api.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		payload = &show

		if prod != show.GUID() {
			return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(messagedef.MsgParameterMismatch, prod, show.GUID()))
		}

		// update the PRODUCTION entry based on resource
		p, err := backend.GetProduction(ctx, show.GUID())
		if err != nil {
			return api.ErrorResponse(c, http.StatusNotFound, err)
		}

		// the attributes we copy from the .yaml
		p.Title = show.Description.Title
		p.Summary = show.Description.Summary
		p.Updated = timestamp.Now()

		if err := backend.UpdateProduction(ctx, p); err != nil {
			return api.ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.EnsureAsset(ctx, show.GUID(), &show.Image); err != nil {
			return api.ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.UpdateShow(ctx, location, show); err != nil {
			return api.ErrorResponse(c, http.StatusBadRequest, err)
		}

	} else if kind == podops.ResourceEpisode {
		var episode *podops.Episode = new(podops.Episode) // FIXME change this !

		if err := c.Bind(episode); err != nil {
			return api.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		payload = &episode

		if prod != episode.Parent() {
			return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(messagedef.MsgParameterMismatch, prod, episode.Parent()))
		}

		// ensure images and media files
		if err := backend.EnsureAsset(ctx, episode.Parent(), &episode.Image); err != nil {
			return api.ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.EnsureAsset(ctx, episode.Parent(), &episode.Enclosure); err != nil {
			return api.ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err := backend.UpdateEpisode(ctx, location, episode); err != nil {
			return api.ErrorResponse(c, http.StatusBadRequest, err)
		}
	} else {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf(messagedef.MsgResourceUnsupportedKind, kind))
	}

	if err := backend.WriteResourceContent(ctx, location, createFlag, forceFlag, payload); err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.Meter(ctx, action, "production", prod, "resource", guid, "kind", kind)

	return api.StandardResponse(c, http.StatusCreated, nil)
}

// DeleteResourceEndpoint deletes a resource and its .yaml file
// GITHUB_ISSUE #14
func DeleteResourceEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c.Request())

	prod := c.Param("prod")
	kind := c.Param("kind")
	guid := c.Param("id")

	/*
		forceFlag := false
		if strings.ToLower(c.QueryParam("f")) == "true" {
			forceFlag = true
		}
	*/
	if !validate.NotEmpty(prod, kind, guid) {
		return api.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}
	if err := AuthorizeAccessResource(ctx, c, ScopeResourceWrite, guid); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	if err := backend.DeleteResource(ctx, prod, kind, guid); err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	// track api access for billing etc
	platform.Meter(ctx, "api.resource.delete", "production", prod, "resource", guid, "kind", kind)

	return c.NoContent(http.StatusNoContent)
}
