package cli

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/podops/podops/podcast"
	"github.com/urfave/cli"
)

const (
	// presetsNameAndPath is the name and location of the config file
	presetsNameAndPath = ".po"

	CmdLineName    = "po"
	CmdLineVersion = "v0.1"

	BasicCmdGroup    = "Basic Commands"
	SettingsCmdGroup = "Settings Commands"
	ShowCmdGroup     = "Show Commands"
	ShowMgmtCmdGroup = "Show Management Commands"

	// All the API & CLI endpoint routes

	// NewShowRoute creates a new production
	NewShowRoute = "/new"
	// CreateRoute creates a resource
	CreateRoute = "/create/:id/:rsrc"
	// UpdateRoute updates a resource
	UpdateRoute = "/update/:id/:rsrc"
)

type (
	// DefaultValues stores all presets the CLI needs
	DefaultValues struct {
		ServiceEndpoint string `json:"url" binding:"required"`
		Token           string `json:"token" binding:"required"`
		ClientID        string `json:"client_id" binding:"required"`
		ShowID          string `json:"show" binding:"required"`
	}
)

var client *podcast.Client

func init() {
	cl, err := podcast.NewClientFromFile(context.Background(), presetsNameAndPath)
	if err != nil {
		log.Fatal(err)
	}
	if cl != nil {
		client = cl
	}
}

// PrintError formats a CLI error and prints it
func PrintError(c *cli.Context, operation string, status int, err error) {
	msg := ""
	switch status {
	case http.StatusInternalServerError:
		msg = fmt.Sprintf("Oops, something went wrong! [%s]", operation)
		break
	case http.StatusConflict:
		msg = fmt.Sprintf("Could not create resource [%s]", operation)
		break
	default:
		msg = err.Error()
	}
	fmt.Println(msg)
}
