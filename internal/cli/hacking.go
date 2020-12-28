package cli

import (
	"context"
	"fmt"

	"github.com/podops/podops/internal/resources"
	"github.com/urfave/cli/v2"
)

// Hack hacks the heck
func Hack(c *cli.Context) error {
	if !client.IsAuthorized() {
		return fmt.Errorf("Not authorized. Use 'po auth' first")
	}

	fmt.Println("Hacking...")
	resources.Build(context.Background(), client.GUID, true)

	return nil
}
