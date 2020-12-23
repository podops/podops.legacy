package cli

import (
	"context"
	"fmt"
	"log"

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
func PrintError(c *cli.Context, err error) {
	msg := fmt.Sprintf("Command '%s'. Something went wrong: %v", c.Command.Name, err)
	fmt.Println(msg)
}
