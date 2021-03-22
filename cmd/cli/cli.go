package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	a "github.com/podops/podops/apiv1"
	cmd "github.com/podops/podops/pkg/cli"
)

const (
	cmdLineName = "po"
)

func main() {

	// initialize CLI
	app := &cli.App{
		Name:    cmdLineName,
		Version: a.VersionString,
		Usage:   fmt.Sprintf("PodOps: Podcast Operations CLI (%s)", a.Version),
		Action: func(c *cli.Context) error {
			fmt.Println(globalHelpText)
			return nil
		},
		Commands: setupCommands(),
		Flags:    globalFlags(),
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setupCommands() []*cli.Command {
	c := []*cli.Command{
		// Basic Commands
		{
			Name:     "list",
			Usage:    "List all productions",
			Category: cmd.BasicCmdGroup,
			Action:   cmd.ListProductionsCommand,
		},
		{
			Name:      "set",
			Usage:     "Sets/shows the default production",
			UsageText: setUsageText,
			Category:  cmd.BasicCmdGroup,
			Action:    cmd.SetProductionsCommand,
		},
		// Settings
		{
			Name:      "login",
			Usage:     "Log in to the service",
			UsageText: "login EMAIL",
			Category:  cmd.SettingsCmdGroup,
			Action:    cmd.LoginCommand,
		},
		{
			Name:     "logout",
			Usage:    "Logout and clear all session information",
			Category: cmd.SettingsCmdGroup,
			Action:   cmd.LogoutCommand,
		},
		{
			Name:      "auth",
			Usage:     "Exchange the token for the API access key",
			UsageText: "auth EMAIL TOKEN",
			Category:  cmd.SettingsCmdGroup,
			Action:    cmd.AuthCommand,
		},
	}
	return c
}

func globalFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:    "prod",
			Usage:   "If present, the podcast scope for the CLI request",
			Aliases: []string{"p"},
		},
	}
	return f
}

//
// all the help texts used in the CLI
//
const (
	globalHelpText = `PodOps: Podcast Operations Client

This client tool helps you to create and produce podcasts.
It also includes administrative commands for managing your live podcasts.

To see the full list of supported commands, run 'po help'`

	setUsageText = `set [GUID]

	 # Display the current show/production
	 po set
	 
	 # Set the current show/production
	 po set GUID`

	getUsageText = `get [RESOURCE]

	 # List all resources
	 po get

	 # List all resources of a type
	 po get [show|episode]

	 # Show details about a resource
	 po get [show|episode] NAME`
)
