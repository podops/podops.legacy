package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/txsvc/commons/pkg/util"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	"github.com/podops/podops/pkg/metadata"
	"github.com/podops/podops/podcast"
)

const (
	// presetsNameAndPath is the name and location of the config file
	presetsNameAndPath = ".po"

	// BasicCmdGroup groups basic commands
	BasicCmdGroup = "\nBasic Commands"
	// SettingsCmdGroup groups settings
	SettingsCmdGroup = "\nSettings Commands"
	// ShowCmdGroup groups basic show commands
	ShowCmdGroup = "\nContent Creation Commands"
	// ShowMgmtCmdGroup groups advanced show commands
	ShowMgmtCmdGroup = "\nContent Management Commands"
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

// remove the local file with login credentials and other state information
func close() error {
	// remove the .po file if it exists
	f, _ := os.Stat(presetsNameAndPath)
	if f != nil {
		err := os.Remove(presetsNameAndPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrintError formats a CLI error and prints it
func PrintError(c *cli.Context, err error) {
	msg := fmt.Sprintf("%s: %v", c.Command.Name, strings.ToLower(err.Error()))
	fmt.Println(msg)
}

// NoOpCommand is just a placeholder
func NoOpCommand(c *cli.Context) error {
	return cli.Exit(fmt.Sprintf("Command '%s' is not implemented", c.Command.Name), 0)
}

// AuthCommand logs into the PodOps service and validates the token
func AuthCommand(c *cli.Context) error {
	token := c.Args().First()

	if token != "" {
		// remove the old settings first
		if err := close(); err != nil {
			return err
		}

		// create a new client and force token verification
		cl, err := podcast.NewClient(context.Background(), token)
		if err != nil {
			fmt.Println("\nNot authorized")
			return nil
		}

		// store the token if valid
		cl.Store(presetsNameAndPath)

		fmt.Println("\nAuthentication successful")
	} else {
		fmt.Println("\nMissing token")
	}

	return nil
}

// LogoutCommand clears all session information
func LogoutCommand(c *cli.Context) error {
	if err := close(); err != nil {
		return err
	}
	client.Close()
	client = nil

	fmt.Println("\nLogout successful")
	return nil
}

// NewProductionCommand requests a new show
func NewProductionCommand(c *cli.Context) error {
	if !client.IsAuthorized() {
		return fmt.Errorf("Not authorized. Use 'po auth' first")
	}

	name := c.Args().First()
	title := c.String("title")
	if title == "" {
		title = "podcast title"
	}
	summary := c.String("summary")
	if summary == "" {
		summary = "podcast summary"
	}

	p, err := client.CreateProduction(name, title, summary)
	if err != nil {
		PrintError(c, err)
		return nil
	}

	show := metadata.DefaultShow(p.Name, title, summary, p.GUID)
	err = dump(fmt.Sprintf("show-%s.yaml", p.GUID), show)
	if err != nil {
		PrintError(c, err)
		return nil
	}

	// update the client
	client.GUID = p.GUID
	client.Store(presetsNameAndPath)

	return nil
}

// ListProductionCommand requests a new show
func ListProductionCommand(c *cli.Context) error {
	if !client.IsAuthorized() {
		return fmt.Errorf("Not authorized. Use 'po auth' first")
	}

	l, err := client.List()
	if err != nil {
		PrintError(c, err)
		return nil
	}

	if len(l.List) == 0 {
		fmt.Println("No shows to list.")
	} else {
		fmt.Println("NAME\t\tGUID\t\tTITLE")
		for _, details := range l.List {
			if details.GUID == client.GUID {
				fmt.Printf("*%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
			} else {
				fmt.Printf("%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
			}
		}
	}

	return nil
}

// SetProductionCommand lists the current show/production, switch to another show/production
func SetProductionCommand(c *cli.Context) error {
	if !client.IsAuthorized() {
		return fmt.Errorf("Not authorized. Use 'po auth' first")
	}

	l, err := client.List()
	if err != nil {
		PrintError(c, err)
		return nil
	}

	if len(l.List) == 0 {
		fmt.Println("No shows available.")
		return nil
	}

	name := c.Args().First()
	if name == "" {
		if client.GUID == "" {
			fmt.Println("No shows selected. Use 'po set NAME' first")
			return nil
		}
		for _, details := range l.List {
			if details.GUID == client.GUID {
				fmt.Println("NAME\t\tGUID\t\tTITLE")
				fmt.Printf("%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
				return nil
			}
		}
		fmt.Println("No shows selected. Use 'po set NAME' first")
		return nil

	}

	for _, details := range l.List {
		if name == details.Name {
			client.GUID = details.GUID
			client.Store(presetsNameAndPath)

			fmt.Println(fmt.Sprintf("Selected '%s'", name))
			fmt.Println("NAME\t\tGUID\t\tTITLE")
			fmt.Printf("%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
			return nil
		}
	}

	fmt.Println(fmt.Sprintf("Can not select '%s'", name))

	return nil
}

// CreateCommand creates a resource from a file, directory or URL
func CreateCommand(c *cli.Context) error {
	if err := client.Valid(); err != nil {
		return err
	}

	if c.NArg() != 1 {
		return fmt.Errorf("Wrong number of arguments. Expected 1, got %d", c.NArg())
	}
	path := c.Args().First()
	force := c.Bool("force")

	resource, kind, guid, err := loadResource(path)
	if err != nil {
		return err
	}

	_, err = client.CreateResource(kind, guid, force, resource)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Created resource %s-%s", kind, guid))
	return nil
}

// UpdateCommand updates a resource from a file, directory or URL
func UpdateCommand(c *cli.Context) error {
	if err := client.Valid(); err != nil {
		return err
	}

	if c.NArg() != 1 {
		return fmt.Errorf("Wrong number of arguments. Expected 1, got %d", c.NArg())
	}
	path := c.Args().First()
	force := c.Bool("force")

	resource, kind, guid, err := loadResource(path)
	if err != nil {
		return err
	}

	_, err = client.UpdateResource(kind, guid, force, resource)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Updated resource %s-%s", kind, guid))
	return nil
}

// TemplateCommand creates a resource template with all default values
func TemplateCommand(c *cli.Context) error {
	template := c.Args().First()
	if template != "show" && template != "episode" {
		fmt.Println(fmt.Sprintf("\nDon't know how to create '%s'", template))
		return nil
	}

	// extract flags or set defaults
	name := c.String("name")
	if name == "" {
		name = "NAME"
	}
	guid := c.String("id")
	if guid == "" {
		guid, _ = util.ShortUUID()
	}
	parent := c.String("parent")
	if parent == "" {
		parent = "PARENT-NAME"
	}
	parentGUID := c.String("parentid")
	if parentGUID == "" {
		parentGUID = "PARENT-ID"
	}

	// create the yamls
	if template == "show" {

		show := metadata.DefaultShow(name, "TITLE", "SUMMARY", guid)
		err := dump(fmt.Sprintf("show-%s.yaml", guid), show)
		if err != nil {
			PrintError(c, err)
			return nil
		}
	} else {

		episode := metadata.DefaultEpisode(name, parent, guid, parentGUID)
		err := dump(fmt.Sprintf("episode-%s.yaml", guid), episode)
		if err != nil {
			PrintError(c, err)
			return nil
		}
	}

	return nil
}

func dump(path string, doc interface{}) error {
	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, data, 0644)
	fmt.Printf("--- %s:\n\n%s\n\n", path, string(data))

	return nil
}
