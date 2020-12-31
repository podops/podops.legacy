package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/podops/podops"
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
		cl, err := podops.NewClient(token)
		if err != nil {
			fmt.Println("\nNot authorized")
			return nil
		}
		err = cl.Store(configName)
		if err != nil {
			fmt.Printf("\nCould not write config. %v\n", err)
			return nil
		}

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
