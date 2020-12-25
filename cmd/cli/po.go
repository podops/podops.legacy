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

func main() {

	// initialize CLI
	app := &cli.App{
		Name:    cmdLineName,
		Version: cmdLineVersion,
		Usage:   "PodOps: Podcast Operations Client",
		Action: func(c *cli.Context) error {
			fmt.Println(globalHelpText)
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
			Flags:     createFlags(),
		},
		{
			Name:      "update",
			Usage:     "Update a resource from a file, directory or URL",
			UsageText: "update FILENAME",
			Category:  cl.ShowCmdGroup,
			Action:    cl.UpdateCommand,
			Flags:     createFlags(),
		},
		{
			Name:     "list",
			Usage:    "List all shows/productions",
			Category: cl.BasicCmdGroup,
			Action:   cl.ListProductionCommand,
		},
		{
			Name:      "set",
			Usage:     "List the current show/production, switch to another show/production",
			UsageText: setUsageText,
			Category:  cl.ShowCmdGroup,
			Action:    cl.SetProductionCommand,
		},
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

func createFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Force create/update",
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

//
// all the help texts used in the CLI
//
const (
	globalHelpText = `PodOps: Podcast Operations Client

This client tool helps you to create and produce podcasts.
It also includes administrative commands for managing your live podcasts.

To see the full list of supported commands, run 'po help'`

	setUsageText = `set [NAME]

	 # Display the current show/production
	 po set
	 
	 # Set the current show/production
	 po set [NAME]`
)
