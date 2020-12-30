package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

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

	// runthe CLI
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
			Usage:    "List all shows/productions",
			Category: cl.BasicCmdGroup,
			Action:   cl.ListProductionCommand,
		},
		{
			Name:      "set",
			Usage:     "List the current show/production, switch to another show/production",
			UsageText: setUsageText,
			Category:  cl.BasicCmdGroup,
			Action:    cl.SetProductionCommand,
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
			Name:      "template",
			Usage:     "Create a resource template with default values",
			UsageText: "template [show|episode]",
			Category:  cl.BasicCmdGroup,
			Action:    cl.TemplateCommand,
			Flags:     templateFlags(),
		},

		// Settings
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
		// Show/production commands
		{
			Name:      "get",
			Usage:     "Lists a single resource/a collection of resources",
			UsageText: getUsageText,
			Category:  cl.ShowCmdGroup,
			Action:    cl.NoOpCommand,
			//Flags:     createFlags(),
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
			Name:      "upload",
			Usage:     "Upload an asset from a file",
			UsageText: "upload FILENAME",
			Category:  cl.ShowCmdGroup,
			Action:    cl.UploadCommand,
			Flags:     createFlags(),
		},
		{
			Name:      "build",
			Usage:     "Start a new build",
			UsageText: "po build",
			Category:  cl.ShowMgmtCmdGroup,
			Action:    cl.BuildCommand,
		},
		{
			Name:      "delete",
			Usage:     "Delete a resource",
			UsageText: "po delete [show|episode] NAME",
			Category:  cl.ShowMgmtCmdGroup,
			Action:    cl.NoOpCommand,
		},
		// HACKING
		{
			Name:     "hack",
			Usage:    "Hacking",
			Category: cl.ShowMgmtCmdGroup,
			Action:   cl.Hack,
		},
	}
	return c
}

func globalFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "Prints debug information",
			Aliases: []string{"d"},
		},
	}
	return f
}

func newShowFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:    "title",
			Usage:   "Show title",
			Aliases: []string{"t"},
		},
		&cli.StringFlag{
			Name:    "summary",
			Usage:   "Show summary",
			Aliases: []string{"s"},
		},
	}
	return f
}

func createFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Usage:   "Force create/update/upload",
			Aliases: []string{"f"},
		},
	}
	return f
}

func templateFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Usage:   "Resource name",
			Aliases: []string{"n"},
		},
		&cli.StringFlag{
			Name:    "parent",
			Usage:   "Parent resource name",
			Aliases: []string{"p"},
		},
		&cli.StringFlag{
			Name:    "guid",
			Usage:   "Resource GUID",
			Aliases: []string{"id"},
		},
		&cli.StringFlag{
			Name:    "parentid",
			Usage:   "Parent resource GUID",
			Aliases: []string{"pid"},
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

	getUsageText = `get [RESOURCE]

	 # List all resources
	 po get

	 # List all resources of a type
	 po get [show|episode]

	 # Show details about a resource
	 po get [show|episode] NAME`
)
