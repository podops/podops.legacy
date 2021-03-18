package auth

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"
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

	DefaultAPIScope = "api.user"

	DefaultTokenType = "user"
)

type (
	// Authorization represents a user, app or bot and its permissions
	Authorization struct {
		ClientID  string `json:"client_id" binding:"required"` // UNIQUE
		Realm     string `json:"realm"`
		Name      string `json:"name"` // DEPREACTED use real instead
		Token     string `json:"token" binding:"required"`
		TokenType string `json:"token_type" binding:"required"` // user,app,bot
		UserID    string `json:"user_id"`                       // depends on TokenType. UserID could equal ClientID or BotUSerID in Slack
		Scope     string `json:"scope"`                         // a comma separated list of scopes, see below
		Expires   int64  `json:"expires"`                       // 0 = never
		// internal
		Revoked bool `json:"-"`
		// FIXME: add revokation flag to the Authorization
		AuthType string `json:"-"` // DEPRECATED currently: jwt, slack
		Created  int64  `json:"-"`
		Updated  int64  `json:"-"`
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
	if a.Revoked == true {
		return false
	}
	if a.Expires < util.Timestamp() {
		return false
	}
	return true
}

// LookupAuthorization looks for an authorization
func LookupAuthorization(ctx context.Context, clientID, realm string) (*Authorization, error) {
	var auth *Authorization
	k := authorizationKey(realm, clientID)

	// FIXME add a cache ?

	if err := platform.DataStore().Get(ctx, k, auth); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil // Not finding one is not an error!
		}
		return nil, err
	}
	return auth, nil
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
