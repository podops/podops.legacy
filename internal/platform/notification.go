package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/fupas/commons/pkg/env"
)

func SendEmail(sender, recipient, subject, body string) error {
	domain := env.GetString("EMAIL_DOMAIN", "m.podops.dev") // FIXME
	if domain == "" {
		return fmt.Errorf("invalid email configuration")
	}
	apiKey := env.GetString("EMAIL_API_KEY", "7424393ba0fbb41e28a4516670af3ddb-203ef6d0-7b16b027") // FIXME
	if apiKey == "" {
		return fmt.Errorf("invalid email configuration")
	}

	mg := mailgun.NewMailgun(domain, apiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)

	message := mg.NewMessage(sender, subject, body, recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)
	if err != nil {
		return err
	}
	fmt.Printf("%s,%s\n", resp, id)
	return nil
}
