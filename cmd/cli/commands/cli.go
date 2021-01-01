package commands

import (
	"fmt"
	"log"
	"os/user"
	"path/filepath"

	"github.com/txsvc/commons/pkg/env"
	"github.com/urfave/cli/v2"

	"github.com/podops/podops"
)

const (
	// BasicCmdGroup groups basic commands
	BasicCmdGroup = "\nBasic Commands"
	// SettingsCmdGroup groups settings
	SettingsCmdGroup = "\nSettings Commands"
	// ShowCmdGroup groups basic show commands
	ShowCmdGroup = "\nContent Creation Commands"
	// ShowMgmtCmdGroup groups advanced show commands
	ShowMgmtCmdGroup = "\nContent Management Commands"

	// configNameAndPath is the name and location of the config file
	configName = "config"
	configPath = ".po"
)

var (
	client             *podops.Client
	defaultPath        string
	defaultPathAndName string
)

func init() {
	path := env.GetString("PODOPS_CREDENTIALS", "")
	if path == "" {
		usr, _ := user.Current()
		defaultPath = filepath.Join(usr.HomeDir, configPath)
		defaultPathAndName = filepath.Join(defaultPath, configName)
	} else {
		defaultPath = filepath.Dir(path)
		defaultPathAndName = path
	}

	cl, err := podops.NewClientFromFile(defaultPathAndName)

	if err != nil {
		log.Fatal(err)
	}
	if cl != nil {
		client = cl
	}
}

// NoOpCommand is just a placeholder
func NoOpCommand(c *cli.Context) error {
	return cli.Exit(fmt.Sprintf("Command '%s' is not implemented", c.Command.Name), 0)
}
