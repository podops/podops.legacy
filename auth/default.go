package auth

import (
	"context"
	"fmt"

	"github.com/txsvc/spa/pkg/env"
	"github.com/txsvc/spa/pkg/timestamp"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/platform"
)

const (
	defaultPodopsScope = "production:read,production:write,production:build,resource:read,resource:write"
)

func CreateSimpleAuthorization(account *Account, req *AuthorizationRequest) *Authorization {
	now := timestamp.Now()

	auth := Authorization{
		ClientID:  account.ClientID,
		Realm:     req.Realm,
		Token:     CreateSimpleToken(),
		TokenType: DefaultTokenType,
		UserID:    req.UserID,
		Scope:     defaultPodopsScope,
		Revoked:   false,
		Expires:   now + (DefaultAuthorizationExpiration * 86400),
		Created:   now,
		Updated:   now,
	}
	return &auth
}

// SendAccountChallenge sends a notification to the user promting to confirm the account
func SendAccountChallenge(ctx context.Context, account *Account) error {
	// FIXME use templates to send a proper email

	url := fmt.Sprintf("%s/login/%s", podops.DefaultAPIEndpoint, account.Ext1)

	if err := platform.SendEmail(env.GetString("EMAIL_FROM", "hello@podops.dev"), account.UserID, "Confirm your account", url); err != nil {
		return err
	}
	return nil
}

// SendAuthToken sends a notification to the user with the current authentication token
func SendAuthToken(ctx context.Context, account *Account) error {
	// FIXME this is not done, just a crude implementation

	if err := platform.SendEmail(env.GetString("EMAIL_FROM", "hello@podops.dev"), account.UserID, "Your confirmation token", account.Ext2); err != nil {
		return err
	}
	return nil
}
