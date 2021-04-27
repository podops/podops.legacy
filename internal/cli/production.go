package cli

import (
	"fmt"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/errordef"
	"github.com/urfave/cli/v2"
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

	show := podops.DefaultShow(p.Name, title, summary, p.GUID, podops.DefaultEndpoint, podops.DefaultCDNEndpoint)
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
		fmt.Println(errordef.MsgCLIProductionsNotFound)
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
		fmt.Println(errordef.MsgCLIProductionsNotFound)
		return nil
	}

	production := c.Args().First()
	if production == "" {
		// print the current show if one has been selected
		if client.DefaultProduction() == "" {
			fmt.Println(errordef.MsgCLIErrorNoProductionSet)
			return nil
		}
		for _, details := range l.Productions {
			if details.GUID == client.DefaultProduction() {
				fmt.Println(productionListing("ID", "NAME", "TITLE", false))
				fmt.Println(productionListing(details.GUID, details.Name, details.Title, false))

				return nil
			}
		}
		fmt.Println(errordef.MsgCLIErrorNoProductionSet)
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

	printMsg(errordef.MsgCLIErrorCanNotSet)
	return nil
}

// BuildCommand starts a new build of the feed
func BuildCommand(c *cli.Context) error {

	prod := getProduction(c)

	build, err := client.Build(prod)
	if err != nil {
		return err
	}

	printMsg(errordef.MsgCLIBuild, prod, build.FeedAliasURL)
	return nil
}
