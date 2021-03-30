package api

import (
	"context"

	"github.com/labstack/echo/v4"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/auth"
	"github.com/podops/podops/pkg/backend"
)

const (
	scopeProductionRead  = "production:read"
	scopeProductionWrite = "production:write"
	scopeProductionBuild = "production:build"
	scopeResourceRead    = "resource:read"
	scopeResourceWrite   = "resource:write"
)

// AuthorizeAccess verifies that the user has the required roles in her authorization
func AuthorizeAccess(ctx context.Context, c echo.Context, scope string) error {
	_, err := auth.CheckAuthorization(ctx, c, scope)
	if err != nil {
		return err
	}

	return nil
}

// AuthorizeAccessProduction verifies that the user has the required roles in
// her authorization and can access the production.
func AuthorizeAccessProduction(ctx context.Context, c echo.Context, scope, claim string) error {
	auth, err := auth.CheckAuthorization(ctx, c, scope)
	if err != nil {
		return err
	}

	p, err := backend.GetProduction(ctx, claim)
	if err != nil {

		return a.ErrNoSuchProduction
	}
	if p.Owner != auth.ClientID {

		return a.ErrNotAuthorized
	}

	return nil
}

// AuthorizeAccessResource verifies that the user has the required roles in
// her authorization and can access the resource.
func AuthorizeAccessResource(ctx context.Context, c echo.Context, scope, claim string) error {
	auth, err := auth.CheckAuthorization(ctx, c, scope)
	if err != nil {
		return err
	}

	r, err := backend.GetResource(ctx, claim)
	if err != nil {
		return a.ErrNoSuchResource
	}
	p, err := backend.GetProduction(ctx, r.ParentGUID)
	if err != nil {
		return a.ErrNoSuchProduction
	}
	if p.Owner != auth.ClientID {
		return a.ErrNotAuthorized
	}

	return nil
}

func validateNotEmpty(claims ...string) bool {
	if claims == nil || len(claims) == 0 {
		return false
	}
	for _, s := range claims {
		if s == "" {
			return false
		}
	}
	return true
}
