package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"

	"github.com/urfave/cli"
)

// AuthCommand logs into the PodOps service and validates the token
func AuthCommand(c *cli.Context) error {
	token := c.Args().First()
	if token != "" {
		// FIXME: validate the token first
		DefaultValuesCLI.Token = token

		StoreConfig()

		fmt.Println("\nAuthentication successful")
	}

	return nil
}

// LogoutCommand clears all session information
func LogoutCommand(c *cli.Context) error {
	df := &DefaultValues{
		ServiceEndpoint: DefaultServiceEndpoint,
		Token:           "",
		ClientID:        "",
		ShowID:          "",
	}
	DefaultValuesCLI = df
	StoreConfig()

	fmt.Println("\nLogout successful")
	return nil
}
