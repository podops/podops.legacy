package auth

import (
	"net/http"

	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/commons/pkg/util"
	"github.com/labstack/echo/v4"
	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/pkg/api"
	"google.golang.org/appengine"
)

// CreateAuthorizationEndpoint creates an authorization and JWT token
func CreateAuthorizationEndpoint(c echo.Context) error {
	var req *a.AuthorizationRequest = new(a.AuthorizationRequest)

	// this endpoint is secured by a master token i.e. a shared secret between
	// the service and the client, NOT a JWT token !!
	bearer := GetBearerToken(c)
	if bearer != env.GetString("MASTER_KEY", "") {
		return c.NoContent(http.StatusUnauthorized)
	}

	err := c.Bind(req)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	token, err := CreateJWTToken(req.Secret, req.Realm, req.ClientID, req.UserID, req.Scope, req.Duration)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	now := util.Timestamp()
	authorization := Authorization{
		ClientID:  req.ClientID,
		Name:      req.Realm,
		Token:     token,
		TokenType: req.ClientType,
		UserID:    req.UserID,
		Scope:     req.Scope,
		Expires:   now + (req.Duration * 86400), // Duration days from now
		AuthType:  AuthTypeJWT,
		Created:   now,
		Updated:   now,
	}
	err = CreateAuthorization(appengine.NewContext(c.Request()), &authorization)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	resp := a.AuthorizationResponse{
		Realm:    req.Realm,
		ClientID: req.ClientID,
		Token:    token,
	}
	return api.StandardResponse(c, http.StatusCreated, &resp)
}

// ValidateAuthorizationEndpoint verifies that the token is valid and exists in the authorization table
func ValidateAuthorizationEndpoint(c echo.Context) error {
	token := GetBearerToken(c)
	if token == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	auth, err := FindAuthorization(appengine.NewContext(c.Request()), token)
	if auth == nil || err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !auth.IsValid() {
		return c.NoContent(http.StatusUnauthorized)
	}
	return c.NoContent(http.StatusAccepted)
}
