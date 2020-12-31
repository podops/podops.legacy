package commands

import (
	"context"
	"fmt"
	"log"
	"os/user"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/podops/podops"
)

const (
	// presetsNameAndPath is the name and location of the config file
	presetsName        = "config"
	presetsNameAndPath = ".po/config"

	// BasicCmdGroup groups basic commands
	BasicCmdGroup = "\nBasic Commands"
	// SettingsCmdGroup groups settings
	SettingsCmdGroup = "\nSettings Commands"
	// ShowCmdGroup groups basic show commands
	ShowCmdGroup = "\nContent Creation Commands"
	// ShowMgmtCmdGroup groups advanced show commands
	ShowMgmtCmdGroup = "\nContent Management Commands"
)

var (
	client *podops.Client
)

func init() {
	usr, _ := user.Current()
	homeDir := filepath.Join(usr.HomeDir, presetsNameAndPath)

	cl, err := podops.NewClientFromFile(context.Background(), homeDir)

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
