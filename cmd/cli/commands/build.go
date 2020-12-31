package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// BuildCommand starts a new build of the feed
func BuildCommand(c *cli.Context) error {

	// FIXME support the 'NAME' option

	url, err := client.Build(client.GUID)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Build production '%s' successful.\nAccess the feed at %s", client.GUID, url))
	return nil
}
