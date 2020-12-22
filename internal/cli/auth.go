package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"context"
	"fmt"
	"os"

	"github.com/podops/podops/podcast"
	"github.com/urfave/cli"
)

// AuthCommand logs into the PodOps service and validates the token
func AuthCommand(c *cli.Context) error {
	token := c.Args().First()
	if token != "" {
		// FIXME: validate the token first
		client.Token = token
		client.Store(presetsNameAndPath)

		fmt.Println("\nAuthentication successful")
	}

	return nil
}

// LogoutCommand clears all session information
func LogoutCommand(c *cli.Context) error {

	err := os.Remove(presetsNameAndPath)
	if err != nil {
		return err
	}

	client, err = podcast.NewClientFromFile(context.Background(), presetsNameAndPath)
	if err != nil {
		return err
	}

	fmt.Println("\nLogout successful")
	return nil
}
