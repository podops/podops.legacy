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
			if details.GUID == client.GUID {
				fmt.Printf("*%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
			} else {
				fmt.Printf("%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
			}
		}
	}

	return nil
}

// SetProductionCommand lists the current show/production, switch to another show/production
func SetProductionCommand(c *cli.Context) error {
	if !client.IsAuthorized() {
		return fmt.Errorf("Not authorized. Use 'po auth' first")
	}

	l, err := client.List()
	if err != nil {
		PrintError(c, err)
		return nil
	}

	if len(l.List) == 0 {
		fmt.Println("No shows available.")
		return nil
	}

	name := c.Args().First()
	if name == "" {
		if client.GUID == "" {
			fmt.Println("No shows selected. Use 'po set NAME' first")
			return nil
		}
		for _, details := range l.List {
			if details.GUID == client.GUID {
				fmt.Println("NAME\t\tGUID\t\tTITLE")
				fmt.Printf("%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
				return nil
			}
		}
		fmt.Println("No shows selected. Use 'po set NAME' first")
		return nil

	}

	for _, details := range l.List {
		if name == details.Name {
			client.GUID = details.GUID
			client.Store(presetsNameAndPath)

			fmt.Println(fmt.Sprintf("Selected '%s'", name))
			fmt.Println("NAME\t\tGUID\t\tTITLE")
			fmt.Printf("%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
			return nil
		}
	}

	fmt.Println(fmt.Sprintf("Can not select '%s'", name))

	return nil
}
