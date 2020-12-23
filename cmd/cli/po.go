package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"

	cl "github.com/podops/podops/internal/cli"
)

const (
	cmdLineName    = "po"
	cmdLineVersion = "v0.1"
)

const helpText = `PodOps: Podcast Operations Client

This client tool helps you to create and produce podcasts.
It also includes administrative commands for managing your live podcasts.

To see the full list of supported commands, run 'po help'`

func main() {

	// initialize CLI
	app := &cli.App{
		Name:    cmdLineName,
		Version: cmdLineVersion,
		Usage:   "PodOps: Podcast Operations Client",
		Action: func(c *cli.Context) error {
			fmt.Println(helpText)
			return nil
		},
		//Flags:    globalFlags(),
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
			Usage:     "Login to the PodOps service",
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
			Usage:     "Create a resource template with default values",
			UsageText: "template [show|episode]",
			Category:  cl.BasicCmdGroup,
			Action:    cl.TemplateCommand,
			Flags:     templateFlags(),
		},
		{
			Name:      "new-show",
			Usage:     "Setup a new show/production",
			UsageText: "new-show NAME",
			Category:  cl.BasicCmdGroup,
			Action:    cl.NewProductionCommand,
			Flags:     newShowFlags(),
		},
		{
			Name:      "create",
			Usage:     "Create a resource from a file, directory or URL",
			UsageText: "create FILENAME",
			Category:  cl.ShowCmdGroup,
			Action:    cl.CreateCommand,
		},
		{
			Name:      "update",
			Usage:     "Update a resource from a file, directory or URL",
			UsageText: "update FILENAME",
			Category:  cl.ShowCmdGroup,
			Action:    cl.UpdateCommand,
		},

		// NOT IMPLEMENTED

		/*
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

		*/
		/*

				{
				Name:     "hack",
				Usage:    "hack",
				Category: cl.BasicCmdGroup,
				Action:   cl.HackCommand,
			},
		*/
	}
	return c
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
			Name:  "id",
			Usage: "Resource GUID",
		},
		&cli.StringFlag{
			Name:  "parentid",
			Usage: "Parent resource GUID",
		},
	}
	return f
}
