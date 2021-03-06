package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fupas/commons/pkg/env"
	"github.com/podops/podops"
	a "github.com/podops/podops/apiv1"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
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

// printError formats a CLI error and prints it
func printError(c *cli.Context, err error) {
	msg := fmt.Sprintf("%s: %v", c.Command.Name, strings.ToLower(err.Error()))
	fmt.Println(msg)
}

func dump(path string, doc interface{}) error {
	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, data, 0644)
	fmt.Printf("--- # %s:\n\n%s\n", path, string(data))

	return nil
}

func loadResource(path string) (interface{}, string, string, error) {
	// FIXME: only local yaml is supported at the moment !

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", "", fmt.Errorf("can not read file '%s': %w", path, err)
	}

	r, kind, guid, err := a.LoadResource(data)
	if err != nil {
		return nil, "", "", err
	}

	return r, kind, guid, nil
}

// removeConfig removes the config file if one exists
func removeConfig() error {
	f, err := os.Stat(defaultPathAndName)
	if err != nil {
		return err
	}
	if f != nil {
		if err := os.Remove(defaultPathAndName); err != nil {
			return err
		}
	}
	return nil
}
