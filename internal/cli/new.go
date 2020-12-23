package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/podops/podops/pkg/metadata"
)

// NewProductionCommand requests a new show
func NewProductionCommand(c *cli.Context) error {
	if err := client.Valid(); err != nil {
		return err
	}

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
		PrintError(c, err)
		return nil
	}

	show := metadata.DefaultShow(p.Name, title, summary, p.GUID)
	err = dump(fmt.Sprintf("show-%s.yaml", p.GUID), show)
	if err != nil {
		PrintError(c, err)
		return nil
	}

	// update the client
	client.GUID = p.GUID
	client.Store(presetsNameAndPath)

	return nil
}
