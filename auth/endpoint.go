package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/pkg/api"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/errordef"
)

// LoginRequestEndpoint initiates the login process.
//
// It creates a new account if the user does not exist and sends
// confirmation request. Once the account is conformed, it will send the
// confirmation token that can be swapped for a real login token.
//
// POST /login
// status 201: new account, account confirmation sent
// status 204: existing account, email with auth token sent
// status 400: invalid request data
// status 403: only logged-out and confirmed users can proceed
func LoginRequestEndpoint(c echo.Context) error {
	var req *AuthorizationRequest = new(AuthorizationRequest)
	ctx := platform.NewHttpContext(c.Request())

	err := c.Bind(req)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if req.Realm == "" || req.UserID == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	account, err := FindAccountByUserID(ctx, req.Realm, req.UserID)
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
		err = ac.accountConfirmNotification(ctx, account)
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
		err = ac.accountConfirmNotification(ctx, account)
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
	err = ac.tokenNotification(ctx, account)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// status 204: existing account, email with token sent
	return c.NoContent(http.StatusNoContent)
}

// LogoutRequestEndpoint removes the session data and state
//
// POST /logout

func LogoutRequestEndpoint(c echo.Context) error {
	var req *AuthorizationRequest = new(AuthorizationRequest)
	ctx := platform.NewHttpContext(c.Request())

	err := c.Bind(req)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if req.Realm == "" || req.UserID == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	token, err := GetBearerToken(c.Request())
	if err != nil {
		return errordef.ErrNoToken
	}
	auth, err := FindAuthorizationByToken(ctx, token)
	if err != nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}
	if auth.UserID != req.UserID || auth.Realm != req.Realm {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	account, err := LookupAccount(ctx, auth.Realm, auth.ClientID)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if account == nil {
		return api.ErrorResponse(c, http.StatusBadRequest, err)
	}

	if account.Status < 0 {
		return c.NoContent(http.StatusForbidden) // account is blocked or deactivated etc ...
	}
	// if we made until here, logout shoud be OK
	account.Status = AccountLoggedOut
	if err := UpdateAccount(ctx, account); err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// LoginConfirmationEndpoint validates an email.
//
// GET /login/:token
// status 307: account is confirmed, redirect to podops.dev/confirmed
// status 400: the request could not be understood by the server due to malformed syntax
// status 401: token is wrong
// status 403: token is expired or has already been used
// status 404: token was not found
func LoginConfirmationEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c.Request())

	token := c.Param("token")
	if token == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}

	account, status, err := ConfirmLoginChallenge(ctx, token)
	if status != http.StatusNoContent {
		return api.ErrorResponse(c, status, err)
	}

	account, err = ResetAuthToken(ctx, account)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	err = ac.tokenNotification(ctx, account)
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// status 307: account is confirmed, email with auth token sent, redirect now
	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/confirmed", podops.DefaultEndpoint))
}

// GetAuthorizationEndpoint exchanges a temporary confirmation token for a 'real' token.
//
// POST /auth
// status 200: success, the real token is in the response
// status 401: token is expired or has already been used, token and user_id do not match
// status 404: token was not found
func GetAuthorizationEndpoint(c echo.Context) error {
	var req *AuthorizationRequest = new(AuthorizationRequest)
	ctx := platform.NewHttpContext(c.Request())

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
