package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"

	"github.com/urfave/cli"
)

type (
	CLINewShowRequest struct {
		Name    string `json:"name" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Summary string `json:"summary" binding:"required"`
	}

	CLINewShowResponse struct {
		GUID string `json:"guid" binding:"required"`
	}
)

// NewShowCommand requests a new show
func NewShowCommand(c *cli.Context) error {
	name := c.Args().First()
	title := c.String("title")
	if title == "" {
		title = "podcast title"
	}
	summary := c.String("summary")
	if summary == "" {
		summary = "podcast summary"
	}

	req := CLINewShowRequest{
		Name:    name,
		Title:   title,
		Summary: summary,
	}
	resp := CLINewShowResponse{}

	err := Post("/new", Token(), &req, &resp)
	if err != nil {
		PrintError(c, err)
	}

	fmt.Println(req)
	fmt.Println(resp)

	return nil
}
