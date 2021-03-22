package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// ListProductionsCommand retrieves all productions or sets the default production
func ListProductionsCommand(c *cli.Context) error {
	l, err := client.Productions()
	if err != nil {
		printError(c, err)
		return nil
	}

	if len(l.Productions) == 0 {
		fmt.Println("No productions found.")
	} else {
		fmt.Println(productionListing("GUID", "NAME", "TITLE", false))
		for _, details := range l.Productions {
			if details.GUID == client.DefaultProduction() {
				fmt.Println(productionListing(details.GUID, details.Name, details.Title, true))
			} else {
				fmt.Println(productionListing(details.GUID, details.Name, details.Title, false))
			}
		}
	}

	return nil
}

// SetProductionsCommand retrieves all productions or sets the default production
func SetProductionsCommand(c *cli.Context) error {

	l, err := client.Productions()
	if err != nil {
		printError(c, err)
		return nil
	}

	if len(l.Productions) == 0 {
		fmt.Println("No shows available.")
		return nil
	}

	production := c.Args().First()
	if production == "" {
		// print the current show if one has been selected
		if client.DefaultProduction() == "" {
			fmt.Println("No production selected. Use 'po set GUID' first")
			return nil
		}
		for _, details := range l.Productions {
			if details.GUID == client.DefaultProduction() {
				fmt.Println(productionListing("GUID", "NAME", "TITLE", false))
				fmt.Println(productionListing(details.GUID, details.Name, details.Title, false))

				return nil
			}
		}
		fmt.Println("No shows selected. Use 'po set NAME' first")
		return nil
	}

	for _, details := range l.Productions {
		if production == details.GUID {
			storeDefaultProduction(production)
			fmt.Println(productionListing("GUID", "NAME", "TITLE", false))
			fmt.Println(productionListing(details.GUID, details.Name, details.Title, true))

			return nil
		}
	}

	fmt.Println(fmt.Sprintf("Can not set production '%s'. Use 'po list' to find available productions", production))

	return nil
}
