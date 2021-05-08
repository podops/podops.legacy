package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/txsvc/platform/v2/pkg/env"

	"github.com/podops/podops"
)

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
