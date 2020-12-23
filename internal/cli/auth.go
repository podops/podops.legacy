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
		if err := close(); err != nil {
			return err
		}

		cl, err := podcast.NewClient(context.Background(), token)
		if err != nil {
			fmt.Println("\nNot authorized")
			return nil
		}

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

func close() error {
	// remove the .po file if it exists
	f, _ := os.Stat(presetsNameAndPath)
	if f != nil {
		err := os.Remove(presetsNameAndPath)
		if err != nil {
			return err
		}
	}
	return nil
}
