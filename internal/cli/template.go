package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"
	"io/ioutil"

	"github.com/txsvc/commons/pkg/util"
	"github.com/urfave/cli/v2"
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

	// extract flags or set defaults
	name := c.String("name")
	if name == "" {
		name = "NAME"
	}
	guid := c.String("id")
	if guid == "" {
		guid, _ = util.ShortUUID()
	}
	parent := c.String("parent")
	if parent == "" {
		parent = "PARENT-NAME"
	}
	parentGUID := c.String("parentid")
	if parentGUID == "" {
		parentGUID = "PARENT-ID"
	}

	// create the yamls
	if template == "show" {

		show := metadata.DefaultShow(name, "TITLE", "SUMMARY", guid)
		err := dump(fmt.Sprintf("show-%s.yaml", guid), show)
		if err != nil {
			PrintError(c, err)
			return nil
		}
	} else {

		episode := metadata.DefaultEpisode(name, parent, guid, parentGUID)
		err := dump(fmt.Sprintf("episode-%s.yaml", guid), episode)
		if err != nil {
			PrintError(c, err)
			return nil
		}
	}

	return nil
}

func dump(path string, doc interface{}) error {
	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, data, 0644)
	fmt.Printf("--- %s:\n\n%s\n\n", path, string(data))

	return nil
}
