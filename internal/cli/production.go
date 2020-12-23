package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/podops/podops/pkg/metadata"
	m "github.com/podops/podops/pkg/metadata"
)

// CreateProductionCommand requests a new show
func CreateProductionCommand(c *cli.Context) error {
	if !client.IsAuthorized() {
		return fmt.Errorf("Not authorized. Use 'po auth' first")
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
	}

	// FIXME replace with /get !!

	show := metadata.DefaultShow(p.Name, title, summary, p.GUID)
	showDoc, err := yaml.Marshal(&show)
	if err != nil {
		PrintError(c, err)
		return nil
	}

	episode := metadata.DefaultEpisode(p.Name, "episode1", p.GUID, p.GUID)
	episodeDoc, err := yaml.Marshal(&episode)
	if err != nil {
		PrintError(c, err)
		return nil
	}

	ioutil.WriteFile(fmt.Sprintf("show-%s.yaml", show.Metadata.Labels[m.LabelGUID]), showDoc, 0644)
	ioutil.WriteFile(fmt.Sprintf("episode-%s.yaml", episode.Metadata.Labels[m.LabelGUID]), episodeDoc, 0644)

	fmt.Printf("--- show dump:\n\n%s\n\n", string(showDoc))
	fmt.Printf("--- episode dump:\n\n%s\n\n", string(episodeDoc))

	// update the client
	client.GUID = p.GUID
	client.Store(presetsNameAndPath)

	return nil
}
