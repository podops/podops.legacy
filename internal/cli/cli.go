package cli

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/podops/podops/podcast"
	"github.com/urfave/cli/v2"
)

const (
	// presetsNameAndPath is the name and location of the config file
	presetsNameAndPath = ".po"

	// BasicCmdGroup groups basic commands
	BasicCmdGroup = "Basic Commands"
	// SettingsCmdGroup groups settings
	SettingsCmdGroup = "Settings Commands"
	// ShowCmdGroup groups basic show commands
	ShowCmdGroup = "Show Commands"
	// ShowMgmtCmdGroup groups advanced show commands
	ShowMgmtCmdGroup = "Show Management Commands"
)

type (
	// ResourceLoaderFunc implements loading of resources
	ResourceLoaderFunc func(data []byte) (interface{}, error)
)

var (
	client *podcast.Client

	resourceLoaders map[string]ResourceLoaderFunc
)

func init() {
	cl, err := podcast.NewClientFromFile(context.Background(), presetsNameAndPath)
	if err != nil {
		log.Fatal(err)
	}
	if cl != nil {
		client = cl
	}

	resourceLoaders = make(map[string]ResourceLoaderFunc)
	resourceLoaders["show"] = loadShowResource
	resourceLoaders["episode"] = loadEpisodeResource
}

// PrintError formats a CLI error and prints it
func PrintError(c *cli.Context, err error) {
	msg := fmt.Sprintf("%s: %v", c.Command.Name, strings.ToLower(err.Error()))
	fmt.Println(msg)
}
