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
			Name:     "shows",
			Usage:    "List all podcasts",
			Category: cmd.BasicCmdGroup,
			Action:   cmd.ListProductionsCommand,
		},
		{
			Name:      "show",
			Usage:     "Sets the default podcast",
			UsageText: setUsageText,
			Category:  cmd.BasicCmdGroup,
			Action:    cmd.SetProductionCommand,
		},
		// resources
		{
			Name:      "create",
			Usage:     "Create a resource from a file, directory or URL",
			UsageText: "create FILENAME",
			Category:  cmd.ShowCmdGroup,
			Action:    cmd.CreateCommand,
			Flags:     createFlags(),
		},
		{
			Name:      "update",
			Usage:     "Update a resource from a file, directory or URL",
			UsageText: "update FILENAME",
			Category:  cmd.ShowCmdGroup,
			Action:    cmd.UpdateCommand,
			Flags:     createFlags(),
		},
		{
			Name:      "get",
			Usage:     "List one or many resources",
			UsageText: getUsageText,
			Category:  cmd.ShowCmdGroup,
			Action:    cmd.GetResourcesCommand,
		},
		{
			Name:      "delete",
			Usage:     "Delete a resource",
			UsageText: "po delete [show|episode] ID",
			Category:  cmd.ShowCmdGroup,
			Action:    cmd.DeleteResourcesCommand,
		},
		{
			Name:      "template",
			Usage:     "Create a resource template with default values",
			UsageText: "template [show|episode] NAME",
			Category:  cmd.ShowCmdGroup,
			Action:    cmd.TemplateCommand,
			Flags:     templateFlags(),
		},
		// build and managment
		{
			Name:      "new",
			Usage:     "Create a new podcast",
			UsageText: "new NAME",
			Category:  cmd.ShowBuildCmdGroup,
			Action:    cmd.NewProductionCommand,
			Flags:     newShowFlags(),
		},
		{
			Name:      "upload",
			Usage:     "Upload an asset from a file",
			UsageText: "upload FILENAME",
			Category:  cmd.ShowBuildCmdGroup,
			Action:    cmd.UploadCommand,
			Flags:     createFlags(),
		},
		{
			Name:      "build",
			Usage:     "Build the podcast feed",
			UsageText: "po build",
			Category:  cmd.ShowBuildCmdGroup,
			Action:    cmd.BuildCommand,
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
			Name:  "guid",
			Usage: "Resource ID",
		},
		&cli.StringFlag{
			Name:  "parent",
			Usage: "Parent resource ID",
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

	setUsageText = `show [ID]

	 # Display the default podcast
	 po show
	 
	 # Sets the default podcast
	 po show ID`

	getUsageText = `list [RESOURCE]

	 # List all resources
	 po list

	 # List all resources of a type
	 po list [shows|episodes|assets]

	 # Show details about a resource
	 po list [show|episode|asset] ID`
)
