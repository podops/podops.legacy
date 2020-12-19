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
	cliName    = "po"
	cliVersion = "v0.1"

	basicCmd    = "Basic Commands"
	settingsCmd = "Settings Commands"
	showCmd     = "Show Commands"
	showMgmtCmd = "Show Management Commands"
)

const helpText = `PodOps: Podcast Operations CLI

This client tool helps you to create and produce podcasts.
It also includes administrative commands for managing your live podcasts.

To see the full list of commands supported, run 'po help'`

var endpoint string = "https://api.podops.dev"

func main() {

	// load presets
	cl.LoadOrCreateDefaultValues()

	// initialize CLI
	app := &cli.App{
		Name:    cliName,
		Version: cliVersion,
		Usage:   "PodOps: Podcast Operations CLI",
		Action: func(c *cli.Context) error {
			fmt.Println(fmt.Sprintf("%s - %s\n", cliName, cliVersion))
			fmt.Println(helpText)
			return nil
		},
		Flags:    setupFlags(),
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

func setupFlags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:        "url",
			Value:       cl.DefaultServiceEndpoint,
			Usage:       "set the service endpoint",
			Destination: &endpoint,
		},
		&cli.StringFlag{
			Name:        "s",
			Usage:       "select the show a command is applied to",
			Destination: &endpoint,
		},
	}
	return f
}

func setupCommands() []cli.Command {
	c := []cli.Command{
		cli.Command{
			Name:     "info",
			Usage:    "Shows an overview of the current show",
			Category: basicCmd,
			Action:   NoopCommand,
		},
		cli.Command{
			Name:      "new-show",
			Usage:     "Request a new show",
			UsageText: "new-show NAME",
			Category:  basicCmd,
			Action:    cl.NewShowCommand,
		},

		cli.Command{
			Name:      "auth",
			Usage:     "Login to the PodOps service and validate the token",
			UsageText: "auth TOKEN",
			Category:  settingsCmd,
			Action:    cl.AuthCommand,
		},
		cli.Command{
			Name:     "logout",
			Usage:    "Logout and clear all session information",
			Category: settingsCmd,
			Action:   cl.LogoutCommand,
		},

		cli.Command{
			Name:     "shows",
			Usage:    "List all shows",
			Category: showCmd,
			Action:   NoopCommand,
		},
		cli.Command{
			Name:     "show",
			Usage:    "Switch to another show",
			Category: showCmd,
			Action:   NoopCommand,
		},

		cli.Command{
			Name:     "create",
			Usage:    "Create a resource from a file, directory or URL",
			Category: showMgmtCmd,
			Action:   NoopCommand,
		},
		cli.Command{
			Name:     "apply",
			Usage:    "Apply a change to a resource from a file, directory or URL",
			Category: showMgmtCmd,
			Action:   NoopCommand,
		},
		cli.Command{
			Name:     "get",
			Usage:    "Display one or many resources by name",
			Category: showMgmtCmd,
			Action:   NoopCommand,
		},
		cli.Command{
			Name:     "delete",
			Usage:    "Delete one or many resources by name",
			Category: showMgmtCmd,
			Action:   NoopCommand,
		},
		cli.Command{
			Name:     "produce",
			Usage:    "Start the production of the podcast feed on the service",
			Category: showMgmtCmd,
			Action:   NoopCommand,
		},
		cli.Command{
			Name:     "build",
			Usage:    "Start the build of the podcast assets locally",
			Category: showMgmtCmd,
			Action:   NoopCommand,
		},
		cli.Command{
			Name:      "template",
			Usage:     "Create a resource template with all default values",
			UsageText: "template show | episode",
			Category:  basicCmd,
			Action:    cl.TemplateCommand,
		},
	}
	return c
}

// NoopCommand does nothing
func NoopCommand(c *cli.Context) error {
	fmt.Println(fmt.Sprintf("%s - %s\n", cliName, cliVersion))
	fmt.Println(fmt.Sprintf("'%s %s' is not yet implemented!", cliName, c.Command.Name))

	return nil
}
