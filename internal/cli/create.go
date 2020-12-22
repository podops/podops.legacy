package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/podops/podops/internal/errors"
	"github.com/podops/podops/pkg/metadata"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// CreateCommand creates a resource from a file, directory or URL
func CreateCommand(c *cli.Context) error {
	if !IsAuthorized() {
		return fmt.Errorf("Not authorized. Use 'po auth' first")
	}

	if c.NArg() != 1 {
		return fmt.Errorf("Wrong number of arguments. Expected 1, got %d", c.NArg())
	}

	resourcePath := strings.ToLower(c.Args().First())

	flagForce := false
	force := strings.ToLower(c.String("force"))
	if force == "yes" || force == "true" {
		flagForce = true
	}

	// FIXME: only local yaml is supported at the moment !
	data, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		return fmt.Errorf("Can not read file '%s'. %w", resourcePath, err)
	}
	var r metadata.BasicResource
	err = yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return fmt.Errorf("Can not read file '%s'. %w", resourcePath, err)
	}

	if r.Kind == "show" {
		var show metadata.Show
		err = yaml.Unmarshal([]byte(data), &show)
		if err != nil {
			return fmt.Errorf("Can not parse file '%s'. %w", resourcePath, err)
		}
		err = show.Validate() // FIXME: only partially implemented !!!
		if err != nil {
			return fmt.Errorf("Resource show is not valid. Reason: %w", err)
		}

		resp := errors.StatusObject{}
		route := fmt.Sprintf("/create/%s/show?force=%v", ShowID(), flagForce)
		status, err := Post(route, Token(), &show, &resp)
		if err != nil {
			PrintError(c, route, status, err)
			return nil
		}

		fmt.Println(fmt.Sprintf("Updated resource '%s'", show.Metadata.Labels[metadata.LabelGUID]))

	} else if r.Kind == "episode" {
		var episode metadata.Episode
		err = yaml.Unmarshal([]byte(data), &episode)
		if err != nil {
			return fmt.Errorf("Can not parse file '%s'. %w", resourcePath, err)
		}
		err = episode.Validate() // FIXME: only partially implemented !!!
		if err != nil {
			return fmt.Errorf("Resource show is not valid. Reason: %w", err)
		}

		resp := errors.StatusObject{}
		route := fmt.Sprintf("/create/%s/episode?force=%v", ShowID(), flagForce)
		status, err := Post(route, Token(), &episode, &resp)
		if err != nil {
			PrintError(c, route, status, err)
			return nil
		}

		fmt.Println(fmt.Sprintf("Updated resource '%s'", episode.Metadata.Labels[metadata.LabelGUID]))

	} else {
		return fmt.Errorf("Unsupported resource type '%s'", r.Kind)
	}

	return nil
}
