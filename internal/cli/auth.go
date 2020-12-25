package cli

import (
	"context"
	"fmt"

	"github.com/podops/podops/podcast"
	"github.com/urfave/cli"
)

// AuthCommand logs into the PodOps service and validates the token
func AuthCommand(c *cli.Context) error {
	token := c.Args().First()

	if token != "" {
		// remove the old settings first
		if err := close(); err != nil {
			return err
		}

		// create a new client and force token verification
		cl, err := podcast.NewClient(context.Background(), token)
		if err != nil {
			fmt.Println("\nNot authorized")
			return nil
		}

		// store the token if valid
		cl.Store(presetsNameAndPath)

		fmt.Println("\nAuthentication successful")
	} else {
		fmt.Println("\nMissing token")
	}

	return nil
}

// LogoutCommand clears all session information
func LogoutCommand(c *cli.Context) error {
	if err := close(); err != nil {
		return err
	}
	client.Close()
	client = nil

	fmt.Println("\nLogout successful")
	return nil
}
