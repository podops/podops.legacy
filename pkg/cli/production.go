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

	show := a.DefaultShow(p.Name, title, summary, p.GUID, a.DefaultPortalEndpoint, a.DefaultCDNEndpoint)
	err = dumpResource(fmt.Sprintf("show-%s.yaml", p.GUID), show)
	if err != nil {
		printError(c, err)
		return nil
	}

	// update the client
	storeDefaultProduction(p.GUID)

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
