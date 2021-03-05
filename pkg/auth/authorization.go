package auth

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"

	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/platform/pkg/platform"
	s "github.com/fupas/platform/pkg/services"
)

const (
	// DatastoreAuthorizations collection AUTHORIZATION
	DatastoreAuthorizations string = "AUTHORIZATIONS"

	// AuthTypeJWT constant jwt
	AuthTypeJWT = "jwt"
	// AuthTypeSlack constant salack
	AuthTypeSlack = "slack"
)

// GetToken returns the oauth token of the workspace integration
func GetToken(ctx context.Context, clientID, authType string) (string, error) {
	// ENV always overrides anything else ...
	token := env.GetString(strings.ToUpper(fmt.Sprintf("%s_AUTH_TOKEN", authType)), "")
	if token != "" {
		return token, nil
	}

	// check the in-memory cache
	key := namedKey(clientID, authType)
	token, _ = s.GetKV(ctx, key)
	if token != "" {
		return token, nil
	}

	auth, err := GetAuthorization(ctx, clientID, authType)
	if err != nil {
		return "", err
	}

	// add the token to the cache
	s.SetKV(ctx, key, auth.Token, 1800)

	return auth.Token, nil
}

// GetAuthorization looks for an authorization. Not finding one is not an error!
func GetAuthorization(ctx context.Context, clientID, authType string) (*Authorization, error) {
	var auth *Authorization
	k := authorizationKey(clientID, authType)

	if err := platform.DataStore().Get(ctx, k, auth); err != nil {
		return nil, err
	}
	return auth, nil
}

// FindAuthorization looks for an authorization by token
func FindAuthorization(ctx context.Context, token string) (*Authorization, error) {
	var auth []*Authorization

	if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreAuthorizations).Filter("Token =", token), &auth); err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, nil
	}
	return auth[0], nil
}

// CreateAuthorization creates all data needed for the OAuth fu
func CreateAuthorization(ctx context.Context, auth *Authorization) error {
	k := authorizationKey(auth.ClientID, auth.AuthType)

	// remove the entry from the cache if it is already there ...
	s.InvalidateKV(ctx, namedKey(auth.ClientID, auth.AuthType))

	// we simply overwrite the existing authorization. If this is no desired, use GetAuthorization first,
	// update the Authorization and then write it back.
	_, err := platform.DataStore().Put(ctx, k, auth)
	return err
}

// authorizationKey creates a datastore key for a workspace authorization based on the team_id.
func authorizationKey(clientID, authType string) *datastore.Key {
	return datastore.NameKey(DatastoreAuthorizations, namedKey(clientID, authType), nil)
}

func namedKey(clientID, authType string) string {
	return authType + "." + clientID
}
