package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/txsvc/platform/v2/pkg/account"
	"github.com/txsvc/platform/v2/pkg/apis/provider"
	auth "github.com/txsvc/platform/v2/pkg/authentication"
	"github.com/txsvc/platform/v2/pkg/env"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/messagedef"
)

const (
	defaultScope = "production:read,production:write,production:build,resource:read,resource:write"
)

type (
	authProviderImpl struct {
	}
)

var (
	PodopsAuthConfig provider.ProviderConfig = provider.WithProvider("platform.podops.auth", provider.TypeAuthentication, PodopsAuthProvider)

	// Interface guards
	_ provider.GenericProvider        = (*authProviderImpl)(nil)
	_ provider.AuthenticationProvider = (*authProviderImpl)(nil)
)

func PodopsAuthProvider() interface{} {
	return &authProviderImpl{}
}

func (a *authProviderImpl) Close() error {
	return nil
}

// AccountChallengeNotification sends a notification to the user promting to confirm the account
func (a *authProviderImpl) AccountChallengeNotification(ctx context.Context, realm, userID string) error {
	acc, err := account.FindAccountByUserID(ctx, realm, userID)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/login/%s", podops.DefaultAPIEndpoint, acc.Token)

	if err := SendEmail(env.GetString("EMAIL_FROM", "hello@podops.dev"), acc.UserID, "Confirm your account", url); err != nil {
		return err
	}
	return nil
}

// ProvideAuthorizationToken sends a notification to the user with the current authentication token
func (a *authProviderImpl) ProvideAuthorizationToken(ctx context.Context, realm, userID, token string) error {
	acc, err := account.FindAccountByUserID(ctx, realm, userID)
	if err != nil {
		return err
	}

	if err := SendEmail(env.GetString("EMAIL_FROM", "hello@podops.dev"), acc.UserID, "Your confirmation token", acc.Token); err != nil {
		return err
	}
	return nil
}

func (a *authProviderImpl) Options() *provider.AuthenticationProviderConfig {
	return &provider.AuthenticationProviderConfig{
		Scope:                    defaultScope,
		Endpoint:                 podops.DefaultEndpoint,
		AuthenticationExpiration: auth.DefaultAuthenticationExpiration,
		AuthorizationExpiration:  auth.DefaultAuthorizationExpiration,
	}
}

func SendEmail(sender, recipient, subject, body string) error {

	if !podops.ValidEmail(recipient) {
		return fmt.Errorf(messagedef.MsgLoginInvalidEmail, recipient)
	}

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
