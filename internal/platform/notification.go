package platform

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/txsvc/spa/pkg/env"
)

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
