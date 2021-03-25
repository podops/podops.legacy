package api

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/pkg/auth"
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
	auth, err := findAuthorization(ctx, c)
	if err != nil {
		return err
	}

	if !hasScope(auth.Scope, scope) {
		return a.ErrNotAuthorized
	}

	return nil
}

// AuthorizeAccessProduction verifies that the user has the required roles in
// her authorization and can access the production.
func AuthorizeAccessProduction(ctx context.Context, c echo.Context, scope, claim string) error {
	auth, err := findAuthorization(ctx, c)
	if err != nil {
		return err
	}

	if !hasScope(auth.Scope, scope) {
		return a.ErrNotAuthorized
	}

	p, err := backend.GetProduction(ctx, claim)
	if err != nil {
		return a.ErrNoSuchProduction
	}
	if p.Owner != auth.UserID {
		return a.ErrNotAuthorized
	}

	return nil
}

// AuthorizeAccessResource verifies that the user has the required roles in
// her authorization and can access the resource.
func AuthorizeAccessResource(ctx context.Context, c echo.Context, scope, claim string) error {
	auth, err := findAuthorization(ctx, c)
	if err != nil {
		return err
	}

	if !hasScope(auth.Scope, scope) {
		return a.ErrNotAuthorized
	}

	r, err := backend.GetResource(ctx, claim)
	if err != nil {
		return a.ErrNoSuchResource
	}
	p, err := backend.GetProduction(ctx, r.ParentGUID)
	if err != nil {
		return a.ErrNoSuchProduction
	}
	if p.Owner != auth.UserID {
		return a.ErrNotAuthorized
	}

	return nil
}

func findAuthorization(ctx context.Context, c echo.Context) (*auth.Authorization, error) {
	token, err := auth.GetBearerToken(c.Request())
	if err != nil {
		return nil, err
	}

	auth, err := auth.FindAuthorizationByToken(ctx, token)
	if err != nil || auth == nil {
		return nil, a.ErrNotAuthorized
	}
	return auth, nil
}

func hasScope(scopes, scope string) bool {
	if scopes == "" || scope == "" {
		return false // empty inputs should never evalute to true
	}

	// FIXME this is a VERY naiv implementation
	return strings.Contains(scopes, scope)
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
