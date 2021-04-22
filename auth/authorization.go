package auth

import (
	"context"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/labstack/echo/v4"

	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"

	"github.com/podops/podops/internal/errordef"
)

const (
	// DatastoreAuthorizations collection AUTHORIZATION
	DatastoreAuthorizations string = "AUTHORIZATIONS"

	// AuthTypeSimpleToken constant token
	AuthTypeSimpleToken = "token"
	// AuthTypeJWT constant jwt
	AuthTypeJWT = "jwt"
	// AuthTypeSlack constant slack
	AuthTypeSlack = "slack"

	// DefaultAuthenticationExpiration in minutes
	DefaultAuthenticationExpiration = 10
	// DefaultAuthorizationExpiration in days
	DefaultAuthorizationExpiration = 90

	// default scopes
	ScopeAPIAdmin = "api:admin"
)

type (
	// Authorization represents a user, app or bot and its permissions
	Authorization struct {
		ClientID  string `json:"client_id" binding:"required"` // UNIQUE
		Realm     string `json:"realm"`
		Token     string `json:"token" binding:"required"`
		TokenType string `json:"token_type" binding:"required"` // user,app,bot
		UserID    string `json:"user_id"`                       // depends on TokenType. UserID could equal ClientID or BotUSerID in Slack
		Scope     string `json:"scope"`                         // a comma separated list of scopes, see below
		Expires   int64  `json:"expires"`                       // 0 = never
		// internal
		Revoked bool  `json:"-"`
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
	}

	// AuthorizationRequest represents a login/authorization request from a user, app, or bot
	AuthorizationRequest struct {
		Realm    string `json:"realm" binding:"required"`
		UserID   string `json:"user_id" binding:"required"`
		ClientID string `json:"client_id"`
		Token    string `json:"token"`
	}
)

func namedKey(part1, part2 string) string {
	return part1 + "." + part2
}

// IsValid verifies that the Authorization is still valid, i.e. is not expired and not revoked.
func (a *Authorization) IsValid() bool {
	if a.Revoked {
		return false
	}
	if a.Expires < util.Timestamp() {
		return false
	}
	return true
}

// HasAdminScope checks if the authorization includes scope 'api:admin'
func (a *Authorization) HasAdminScope() bool {
	return strings.Contains(a.Scope, ScopeAPIAdmin)
}

// GetBearerToken extracts the bearer token
func GetBearerToken(r *http.Request) (string, error) {

	// FIXME optimize this !!

	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return "", errordef.ErrNoToken
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 {
		return "", errordef.ErrNoToken
	}
	if parts[0] == "Bearer" {
		return parts[1], nil
	}

	return "", errordef.ErrNoToken
}

// GetClientID extracts the ClientID from the token
func GetClientID(ctx context.Context, r *http.Request) (string, error) {
	token, err := GetBearerToken(r)
	if err != nil {
		return "", err
	}

	// FIXME optimize this, e.g. implement caching

	auth, err := FindAuthorizationByToken(ctx, token)
	if err != nil {
		return "", err
	}
	if auth == nil {
		return "", errordef.ErrNotAuthorized
	}

	return auth.ClientID, nil
}

// CheckAuthorization relies on the presence of a bearer token and validates the
// matching authorization against a list of requested scopes. If everything checks
// out, the function returns the authorization or an error otherwise.
func CheckAuthorization(ctx context.Context, c echo.Context, scope string) (*Authorization, error) {
	token, err := GetBearerToken(c.Request())
	if err != nil {
		return nil, err
	}

	auth, err := FindAuthorizationByToken(ctx, token)
	if err != nil || auth == nil || !auth.IsValid() {
		return nil, errordef.ErrNotAuthorized
	}

	account, err := FindAccountByUserID(ctx, auth.Realm, auth.UserID)
	if err != nil {
		return nil, err
	}

	if account.Status != AccountActive {
		return nil, errordef.ErrNotAuthorized // not logged-in
	}

	if !hasScope(auth.Scope, scope) {
		return nil, errordef.ErrNotAuthorized
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

// LookupAuthorization looks for an authorization
func LookupAuthorization(ctx context.Context, realm, clientID string) (*Authorization, error) {
	var auth Authorization
	k := authorizationKey(realm, clientID)

	// FIXME add a cache ?

	if err := platform.DataStore().Get(ctx, k, &auth); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil // Not finding one is not an error!
		}
		return nil, err
	}
	return &auth, nil
}

// FindAuthorizationByToken looks for an authorization by the token
func FindAuthorizationByToken(ctx context.Context, token string) (*Authorization, error) {
	var auth []*Authorization

	// FIXME add a cache ?

	if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreAuthorizations).Filter("Token =", token), &auth); err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, nil
	}
	return auth[0], nil
}

// CreateAuthorization creates all data needed for the auth fu
func CreateAuthorization(ctx context.Context, auth *Authorization) error {
	k := authorizationKey(auth.Realm, auth.ClientID)

	// FIXME add a cache ?

	// we simply overwrite the existing authorization. If this is no desired, use GetAuthorization first,
	// update the Authorization and then write it back.
	_, err := platform.DataStore().Put(ctx, k, auth)
	return err
}

// UpdateAuthorization updates all data needed for the auth fu
func UpdateAuthorization(ctx context.Context, auth *Authorization) error {
	k := authorizationKey(auth.Realm, auth.ClientID)
	// FIXME add a cache ?

	// we simply overwrite the existing authorization. If this is no desired, use GetAuthorization first,
	// update the Authorization and then write it back.
	_, err := platform.DataStore().Put(ctx, k, auth)
	return err
}

// authorizationKey creates a datastore key for a workspace authorization based on the team_id.
func authorizationKey(realm, client string) *datastore.Key {
	return datastore.NameKey(DatastoreAuthorizations, namedKey(realm, client), nil)
}
