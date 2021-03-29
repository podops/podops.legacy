package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"

	a "github.com/podops/podops/apiv1"
)

// NewProductionCommand creates a new podcast
func NewProductionCommand(c *cli.Context) error {

	name := c.Args().First()
	title := c.String("title")
	if title == "" {
		title = "podcast title"
	}
	summary := c.String("summary")
	if summary == "" {
		summary = "podcast summary"
	}

	p, err := client.CreateProduction(name, title, summary)
	if err != nil {
		printError(c, err)
		return nil
	}

	show := DefaultShow(p.Name, title, summary, p.GUID, a.DefaultEndpoint, a.DefaultCDNEndpoint)
	err = dumpResource(fmt.Sprintf("show-%s.yaml", p.GUID), show)
	if err != nil {
		printError(c, err)
		return nil
	}

	// update the client
	storeDefaultProduction(p.GUID)

	return nil
}

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
		fmt.Println(productionListing("ID", "NAME", "TITLE", false))
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

// SetProductionCommand retrieves all productions or sets the default production
func SetProductionCommand(c *cli.Context) error {

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
			fmt.Println("No production set. Use 'po set ID' first")
			return nil
		}
		for _, details := range l.Productions {
			if details.GUID == client.DefaultProduction() {
				fmt.Println(productionListing("ID", "NAME", "TITLE", false))
				fmt.Println(productionListing(details.GUID, details.Name, details.Title, false))

				return nil
			}
		}
		fmt.Println("No production set. Use 'po set ID' first")
		return nil
	}

	for _, details := range l.Productions {
		if production == details.GUID {
			storeDefaultProduction(production)
			fmt.Println(productionListing("ID", "NAME", "TITLE", false))
			fmt.Println(productionListing(details.GUID, details.Name, details.Title, true))

			return nil
		}
	}

	fmt.Println(fmt.Sprintf("Can not set production '%s'. Use 'po list' to find available productions", production))

	return nil
}

// BuildCommand starts a new build of the feed
func BuildCommand(c *cli.Context) error {

	prod := getProduction(c)

	build, err := client.Build(prod)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Build production '%s' successful.\nAccess the feed at %s", prod, build.FeedAliasURL))
	return nil
}
