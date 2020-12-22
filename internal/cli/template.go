package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/txsvc/commons/pkg/util"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/podops/podops/pkg/metadata"
)

// TemplateCommand creates a resource template with all default values
func TemplateCommand(c *cli.Context) error {
	template := c.Args().First()
	if template != "show" && template != "episode" {
		fmt.Println(fmt.Sprintf("\nDon't know how to create '%s'", template))
		return nil
	}

	if template == "show" {

		guid, _ := util.ShortUUID()
		show := metadata.DefaultShow("podcast-name", "Podcast Title", "Podcast summary describing the podcast", guid)
		showDoc, err := yaml.Marshal(&show)
		if err != nil {
			PrintError(c, NewShowRoute, http.StatusInternalServerError, err)
			return nil
		}

		ioutil.WriteFile(fmt.Sprintf("show-%s.yaml", guid), showDoc, 0644)
		fmt.Printf("--- show dump:\n\n%s\n\n", string(showDoc))
	} else {

		guid, _ := util.ShortUUID()
		episode := metadata.DefaultEpisode("podcast-name", "episode1", guid, "parent-"+guid)
		episodeDoc, err := yaml.Marshal(&episode)
		if err != nil {
			PrintError(c, NewShowRoute, http.StatusInternalServerError, err)
			return nil
		}

		ioutil.WriteFile(fmt.Sprintf("episode-%s.yaml", guid), episodeDoc, 0644)
		fmt.Printf("--- episode dump:\n\n%s\n\n", string(episodeDoc))
	}

	return nil
}
