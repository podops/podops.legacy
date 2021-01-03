package commands

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	a "github.com/podops/podops/apiv1"
)

// NewProductionCommand requests a new show
func NewProductionCommand(c *cli.Context) error {

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
		printError(c, err)
		return nil
	}

	show := a.DefaultShow(client.ServiceEndpoint, p.Name, title, summary, p.GUID)
	err = dump(fmt.Sprintf("show-%s.yaml", p.GUID), show)
	if err != nil {
		printError(c, err)
		return nil
	}

	// update the client
	client.GUID = p.GUID
	client.Store(defaultPathAndName)

	return nil
}

// ListProductionsCommand retrieves all shows
func ListProductionsCommand(c *cli.Context) error {

	l, err := client.Productions()
	if err != nil {
		printError(c, err)
		return nil
	}

	if len(l.Productions) == 0 {
		fmt.Println("No shows to list.")
	} else {
		fmt.Println("NAME\t\tGUID\t\tTITLE")
		for _, details := range l.Productions {
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

	l, err := client.Productions()
	if err != nil {
		printError(c, err)
		return nil
	}

	if len(l.Productions) == 0 {
		fmt.Println("No shows available.")
		return nil
	}

	name := c.Args().First()
	if name == "" {
		if client.GUID == "" {
			fmt.Println("No shows selected. Use 'po set NAME' first")
			return nil
		}
		for _, details := range l.Productions {
			if details.GUID == client.GUID {
				fmt.Println("NAME\t\tGUID\t\tTITLE")
				fmt.Printf("%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
				return nil
			}
		}
		fmt.Println("No shows selected. Use 'po set NAME' first")
		return nil

	}

	for _, details := range l.Productions {
		if name == details.Name {
			client.GUID = details.GUID
			client.Store(defaultPathAndName)

			fmt.Println(fmt.Sprintf("Selected '%s'", name))
			fmt.Println("NAME\t\tGUID\t\tTITLE")
			fmt.Printf("%s\t\t%s\t%s\n", details.Name, details.GUID, details.Title)
			return nil
		}
	}

	fmt.Println(fmt.Sprintf("Can not select '%s'", name))

	return nil
}

// ListResourcesCommand requests a new show
func ListResourcesCommand(c *cli.Context) error {

	kind := strings.ToLower(c.Args().First())

	if c.NArg() < 2 {
		// get a list of resources
		l, err := client.Resources(client.GUID, kind)
		if err != nil {
			printError(c, err)
			return nil
		}

		if len(l.Resources) == 0 {
			fmt.Println("No resources to list.")
		} else {
			fmt.Println("NAME\t\tGUID\t\tKIND")
			for _, details := range l.Resources {
				fmt.Printf("%s\t\t%s\t%s\n", details.Name, details.GUID, details.Kind)
			}
		}
	} else {
		// get a single resource
		guid := c.Args().Get(1)

		var rsrc interface{}
		err := client.Resource(client.GUID, kind, guid, &rsrc)
		if err != nil {
			printError(c, err)
			return nil
		}

		// FIXME verify that rsrc.Kind == kind

		data, err := yaml.Marshal(rsrc)
		if err != nil {
			return err
		}

		fmt.Printf("\n--- %s/%s-%s:\n\n%s\n\n", client.GUID, kind, guid, string(data))

	}

	return nil
}

// DeleteResourcesCommand deletes a resource
func DeleteResourcesCommand(c *cli.Context) error {

	if c.NArg() != 2 {
		return fmt.Errorf("Wrong number of arguments. Expected 2, got %d", c.NArg())
	}

	kind := strings.ToLower(c.Args().First())
	guid := c.Args().Get(1)

	status, err := client.Delete(client.GUID, kind, guid)
	if status > http.StatusAccepted && err == nil {
		fmt.Println(fmt.Sprintf("Could not delete resource '%s/%s-%s'", client.GUID, kind, guid))
		return nil
	}
	if err != nil {
		printError(c, err)
		return err
	}

	fmt.Println(fmt.Sprintf("Successfully delete resource '%s/%s-%s'", client.GUID, kind, guid))
	return nil
}
