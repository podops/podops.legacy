package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"

	cl "github.com/podops/podops/internal/cli"
)

const helpText = `PodOps: Podcast Operations CLI

This client tool helps you to create and produce podcasts.
It also includes administrative commands for managing your live podcasts.

To see the full list of commands supported, run 'po help'`

var endpoint string = "https://api.podops.dev"

func main() {

	// load presets
	cl.LoadOrCreateConfig()

	// initialize CLI
	app := &cli.App{
		Name:    cl.CmdLineName,
		Version: cl.CmdLineVersion,
		Usage:   "PodOps: Podcast Operations CLI",
		Action: func(c *cli.Context) error {
			fmt.Println(fmt.Sprintf("%s - %s\n", cl.CmdLineName, cl.CmdLineVersion))
			fmt.Println(helpText)
			return nil
		},
		Flags:    globalFlags(),
		Commands: setupCommands(),
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	// runthe CLI
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setupCommands() []cli.Command {
	c := []cli.Command{
		{
			Name:      "auth",
			Usage:     "Login to the PodOps service and validate the token",
			UsageText: "auth TOKEN",
			Category:  cl.SettingsCmdGroup,
			Action:    cl.AuthCommand,
		},
		{
			Name:     "logout",
			Usage:    "Logout and clear all session information",
			Category: cl.SettingsCmdGroup,
			Action:   cl.LogoutCommand,
		},
		{
			Name:      "template",
			Usage:     "Create a resource template with all default values",
			UsageText: "template [show | episode]",
			Category:  cl.BasicCmdGroup,
			Action:    cl.TemplateCommand,
			Flags:     templateFlags(),
		},
		{
			Name:      "new-show",
			Usage:     "Setup a new show",
			UsageText: "new-show NAME",
			Category:  cl.BasicCmdGroup,
			Action:    cl.NewShowCommand,
			Flags:     newShowFlags(),
		},
		{
			Name:     "create",
			Usage:    "Create a resource from a file, directory or URL",
			Category: cl.ShowCmdGroup,
			Action:   cl.CreateCommand,
			Flags:    createFlags(),
		},
		// NOT IMPLEMENTED
		{
			Name:     "info",
			Usage:    "Shows an overview of the current show",
			Category: cl.BasicCmdGroup,
			Action:   cl.NoopCommand,
		},

		{
			Name:     "shows",
			Usage:    "List all shows",
			Category: cl.ShowCmdGroup,
			Action:   cl.NoopCommand,
		},
		{
			Name:     "show",
			Usage:    "Switch to another show",
			Category: cl.ShowCmdGroup,
			Action:   cl.NoopCommand,
		},
		{
			Name:     "apply",
			Usage:    "Apply a change to a resource from a file, directory or URL",
			Category: cl.ShowMgmtCmdGroup,
			Action:   cl.NoopCommand,
		},
		{
			Name:     "get",
			Usage:    "Display one or many resources by name",
			Category: cl.ShowMgmtCmdGroup,
			Action:   cl.NoopCommand,
		},
		{
			Name:     "delete",
			Usage:    "Delete one or many resources by name",
			Category: cl.ShowMgmtCmdGroup,
			Action:   cl.NoopCommand,
		},
		{
			Name:     "produce",
			Usage:    "Start the production of the podcast feed on the service",
			Category: cl.ShowMgmtCmdGroup,
			Action:   cl.NoopCommand,
		},
		{
			Name:     "build",
			Usage:    "Start the build of the podcast assets locally",
			Category: cl.ShowMgmtCmdGroup,
			Action:   cl.NoopCommand,
		},
	}
	return c
}

func globalFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:        "url",
			Value:       cl.DefaultServiceEndpoint,
			Usage:       "set the service endpoint",
			Destination: &cl.DefaultValuesCLI.ServiceEndpoint,
		},
	}
	return f
}

func newShowFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:  "title",
			Usage: "Show title",
		},
		&cli.StringFlag{
			Name:  "summary",
			Usage: "Show summary",
		},
	}
	return f
}

func templateFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "Resource name",
		},
		&cli.StringFlag{
			Name:  "parent",
			Usage: "Parent resource name",
		},
		&cli.StringFlag{
			Name:  "pid",
			Usage: "Parent resource GUID",
		},
	}
	return f
}

func createFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:     "force",
			Usage:    "[yes | no] Forces changes to a resource",
			Required: false,
			Value:    "no",
		},
	}
	return f
}
