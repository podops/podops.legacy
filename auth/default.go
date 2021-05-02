package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/txsvc/platform/v2/pkg/env"
	"github.com/txsvc/platform/v2/pkg/timestamp"

	"github.com/podops/podops"
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

	if err := SendEmail(env.GetString("EMAIL_FROM", "hello@podops.dev"), account.UserID, "Confirm your account", url); err != nil {
		return err
	}
	return nil
}

// SendAuthToken sends a notification to the user with the current authentication token
func SendAuthToken(ctx context.Context, account *Account) error {
	// FIXME this is not done, just a crude implementation

	if err := SendEmail(env.GetString("EMAIL_FROM", "hello@podops.dev"), account.UserID, "Your confirmation token", account.Ext2); err != nil {
		return err
	}
	return nil
}

func SendEmail(sender, recipient, subject, body string) error {
	domain := env.GetString("EMAIL_DOMAIN", "")
	apiKey := env.GetString("EMAIL_API_KEY", "")

	mg := mailgun.NewMailgun(domain, apiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)

	message := mg.NewMessage(sender, subject, body, recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
