package apiv1

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/podops/podops/auth"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/errordef"
)

const (
	ScopeProductionRead  = "production:read"
	ScopeProductionWrite = "production:write"
	ScopeProductionBuild = "production:build"
	ScopeResourceRead    = "resource:read"
	ScopeResourceWrite   = "resource:write"
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

	if auth.HasAdminScope() {
		// can access any production
		return nil
	}

	p, err := backend.GetProduction(ctx, claim)
	if err != nil || p == nil {
		return errordef.ErrNoSuchProduction
	}

	if p.Owner != auth.ClientID {
		return errordef.ErrNotAuthorized
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

	if auth.HasAdminScope() {
		// can access any resource
		return nil
	}

	r, err := backend.GetResource(ctx, claim)
	if err != nil || r == nil {
		return errordef.ErrNoSuchResource
	}

	p, err := backend.GetProduction(ctx, r.ParentGUID)
	if err != nil || p == nil {
		return errordef.ErrNoSuchProduction
	}
	if p.Owner != auth.ClientID {
		return errordef.ErrNotAuthorized
	}

	return nil
}

// ValidateNotEmpty test all provided values to be not empty
func ValidateNotEmpty(claims ...string) bool {
	if len(claims) == 0 {
		return false
	}
	for _, s := range claims {
		if s == "" {
			return false
		}
	}
	return true
}
