package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/authentication"
	"github.com/txsvc/platform/v2/pkg/account"
	"github.com/txsvc/platform/v2/pkg/env"

	"github.com/podops/podops"
)

const (
	defaultScope = "production:read,production:write,production:build,resource:read,resource:write"
)

type (
	authProviderImpl struct {
	}
)

var (
	PodopsAuthConfig platform.PlatformOpts = platform.WithProvider("platform.podops.auth", platform.ProviderTypeAuthentication, PodopsAuthProvider)

	// Interface guards
	_ platform.GenericProvider              = (*authProviderImpl)(nil)
	_ authentication.AuthenticationProvider = (*authProviderImpl)(nil)
)

func PodopsAuthProvider() interface{} {
	return &authProviderImpl{}
}

func (a *authProviderImpl) Close() error {
	return nil
}

// AccountChallengeNotification sends a notification to the user promting to confirm the account
func (a *authProviderImpl) AccountChallengeNotification(ctx context.Context, account *account.Account) error {
	url := fmt.Sprintf("%s/login/%s", podops.DefaultAPIEndpoint, account.Token)

	if err := SendEmail(env.GetString("EMAIL_FROM", "hello@podops.dev"), account.UserID, "Confirm your account", url); err != nil {
		return err
	}
	return nil
}

// ProvideAuthorizationToken sends a notification to the user with the current authentication token
func (a *authProviderImpl) ProvideAuthorizationToken(ctx context.Context, account *account.Account) error {
	if err := SendEmail(env.GetString("EMAIL_FROM", "hello@podops.dev"), account.UserID, "Your confirmation token", account.Token); err != nil {
		return err
	}
	return nil
}

func (a *authProviderImpl) Options() *authentication.AuthenticationProviderOpts {
	return &authentication.AuthenticationProviderOpts{
		Scope:                    defaultScope,
		Endpoint:                 podops.DefaultEndpoint,
		AuthenticationExpiration: authentication.DefaultAuthenticationExpiration,
		AuthorizationExpiration:  authentication.DefaultAuthorizationExpiration,
	}
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
