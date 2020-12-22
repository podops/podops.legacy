package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/podops/podops/pkg/metadata"
	m "github.com/podops/podops/pkg/metadata"
)

type (
	// NewShowRequest defines the request
	NewShowRequest struct {
		Name    string `json:"name" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Summary string `json:"summary" binding:"required"`
	}

	// NewShowResponse defines the request
	NewShowResponse struct {
		Name string `json:"name" binding:"required"`
		GUID string `json:"guid" binding:"required"`
	}
)

// NewShowCommand requests a new show
func NewShowCommand(c *cli.Context) error {
	if !IsAuthorized() {
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

	req := NewShowRequest{
		Name:    name,
		Title:   title,
		Summary: summary,
	}

	resp := NewShowResponse{}
	status, err := Post(NewShowRoute, Token(), &req, &resp)
	if err != nil {
		PrintError(c, NewShowRoute, status, err)
		return nil
	}

	show := metadata.DefaultShow(resp.Name, title, summary, resp.GUID)
	showDoc, err := yaml.Marshal(&show)
	if err != nil {
		PrintError(c, NewShowRoute, http.StatusInternalServerError, err)
		return nil
	}

	episode := metadata.DefaultEpisode(resp.Name, "episode1", resp.GUID, resp.GUID)
	episodeDoc, err := yaml.Marshal(&episode)
	if err != nil {
		PrintError(c, NewShowRoute, http.StatusInternalServerError, err)
		return nil
	}

	ioutil.WriteFile(fmt.Sprintf("show-%s.yaml", show.Metadata.Labels[m.LabelGUID]), showDoc, 0644)
	ioutil.WriteFile(fmt.Sprintf("episode-%s.yaml", episode.Metadata.Labels[m.LabelGUID]), episodeDoc, 0644)

	fmt.Printf("--- show dump:\n\n%s\n\n", string(showDoc))
	fmt.Printf("--- episode dump:\n\n%s\n\n", string(episodeDoc))

	// update the config
	DefaultValuesCLI.DefaultShow = resp.GUID

	return nil
}
