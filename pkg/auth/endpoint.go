package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/appengine"

	"github.com/podops/podops/pkg/api"
)

// LoginEndpoint initiates the login process.
//
// It creates a new account if the user does not exist and sends
// confirmation request. Once the account is conformed, it will send the
// confirmation token that can be swapped for a real login token.
//
// POST /login
// status 201: new account, account confirmation sent
// status 204: existing account, email with auth token sent
// status 403: only logged-out and confirmed users can proceed
func LoginEndpoint(c echo.Context) error {
	var req *AuthorizationRequest = new(AuthorizationRequest)
	ctx := appengine.NewContext(c.Request())

	err := c.Bind(req)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if req.Realm == "" || req.UserID == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	account, err := LookupAccount(ctx, req.Realm, req.UserID)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// new account
	if account == nil {
		// #1: create a new account
		account, err = CreateAccount(ctx, req.Realm, req.UserID)
		if err != nil {
			return api.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		// #2: send the confirmation link
		err = SendAccountChallenge(ctx, account)
		if err != nil {
			return api.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		// status 201: new account
		return c.NoContent(http.StatusCreated)
	}

	// existing account but check some stuff first ...
	if account.Confirmed == 0 {
		// #1: update the expiration timestamp
		account, err = ResetAccountChallenge(ctx, account)
		if err != nil {
			return api.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		// #2: send the account confirmation link
		err = SendAccountChallenge(ctx, account)
		if err != nil {
			return api.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		// status 201: new account
		return c.NoContent(http.StatusCreated)
	}
	if account.Status != 0 {
		// status 403: only logged-out and confirmed users can proceed, do nothing otherwise
		return api.ErrorResponse(c, http.StatusForbidden, err)
	}

	// create and send the auth token
	account, err = ResetAuthToken(ctx, account)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	err = SendAuthToken(ctx, account)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// status 204: existing account, email with token sent
	return c.NoContent(http.StatusNoContent)
}

// LoginConfirmationEndpoint validates an email.
//
// GET /login/:token
// status 204: account is confirmed, next step started
// status 400: the request could not be understood by the server due to malformed syntax
// status 401: token is wrong
// status 403: token is expired or has already been used
// status 404: token was not found
func LoginConfirmationEndpoint(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())

	token := c.Param("token")
	if token == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':token"))
	}

	account, status, err := ConfirmLoginChallenge(ctx, token)
	if status != http.StatusNoContent {
		return api.ErrorResponse(c, status, err)
	}

	account, err = ResetAuthToken(ctx, account)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	err = SendAuthToken(ctx, account)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// status 204: account is confirmed, email with auth token sent
	return c.NoContent(http.StatusNoContent)
}

// GetAuthorizationEndpoint exchanges a temporary confirmation token for a 'real' token.
//
// POST /auth
// status 200: success, the real token is in the response
// status 401: token is expired or has already been used, token and user_id do not match
// status 404: token was not found
func GetAuthorizationEndpoint(c echo.Context) error {
	var req *AuthorizationRequest = new(AuthorizationRequest)
	ctx := appengine.NewContext(c.Request())

	err := c.Bind(req)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if req.Token == "" || req.Realm == "" || req.UserID == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	auth, status, err := exchangeToken(ctx, req, c.Request().RemoteAddr)
	if status != http.StatusOK {
		return api.ErrorResponse(c, status, err)
	}

	req.Token = auth.Token
	req.ClientID = auth.ClientID

	return api.StandardResponse(c, status, req)
}

// CreateAuthorizationEndpoint creates an authorization and JWT token
func CreateAuthorizationEndpoint(c echo.Context) error {
	return nil
}
