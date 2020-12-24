package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"

	"github.com/urfave/cli"
)

// ListProductionCommand requests a new show
func ListProductionCommand(c *cli.Context) error {
	if !client.IsAuthorized() {
		return fmt.Errorf("Not authorized. Use 'po auth' first")
	}

	l, err := client.List()
	if err != nil {
		PrintError(c, err)
		return nil
	}

	if len(l.List) == 0 {
		fmt.Println("No shows to list.")
	} else {
		fmt.Println("NAME\t\tGUID\t\tTITLE")
		for _, details := range l.List {
			fmt.Printf("%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
		}
	}

	return nil
}
