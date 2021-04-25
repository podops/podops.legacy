package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/fupas/commons/pkg/env"
)

func SendEmail(sender, recipient, subject, body string) error {
	domain := env.GetString("EMAIL_DOMAIN", "") // FIXME. Check this on startup and not here.
	if domain == "" {
		return fmt.Errorf("invalid email configuration")
	}
	apiKey := env.GetString("EMAIL_API_KEY", "") // FIXME. Check this on startup and not here.
	if apiKey == "" {
		return fmt.Errorf("invalid email configuration")
	}

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
