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

	show := a.DefaultShow(p.Name, title, summary, p.GUID, a.DefaultPortalEndpoint, a.DefaultCDNEndpoint)
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
		fmt.Println(productionListing("GUID", "NAME", "TITLE", false))
		for _, details := range l.Productions {
			if details.GUID == client.GUID {
				fmt.Println(productionListing(details.GUID, details.Name, details.Title, true))
			} else {
				fmt.Println(productionListing(details.GUID, details.Name, details.Title, false))
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
		// print the current show if one has been selected
		if client.GUID == "" {
			fmt.Println("No shows selected. Use 'po set NAME' first")
			return nil
		}
		for _, details := range l.Productions {
			if details.GUID == client.GUID {
				fmt.Println(productionListing("GUID", "NAME", "TITLE", false))
				fmt.Println(productionListing(details.GUID, details.Name, details.Title, false))

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

			fmt.Println(productionListing("GUID", "NAME", "TITLE", false))
			fmt.Println(productionListing(details.GUID, details.Name, details.Title, true))

			return nil
		}
	}

	fmt.Println(fmt.Sprintf("Can not set show '%s'. Use 'po list' to list available shows", name))

	return nil
}

// ListResourcesCommand list all resource associated with a show
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
			fmt.Println(assetListing("GUID", "NAME", "KIND"))
			for _, details := range l.Resources {
				fmt.Println(assetListing(details.GUID, details.Name, details.Kind))
			}
		}
	} else {
		// get a single resource
		guid := c.Args().Get(1)

		var rsrc interface{}
		err := client.GetResource(client.GUID, kind, guid, &rsrc)
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
		return fmt.Errorf("wrong number of arguments: expected 2, got %d", c.NArg())
	}

	kind := strings.ToLower(c.Args().First())
	guid := c.Args().Get(1)

	status, err := client.DeleteResource(client.GUID, kind, guid)
	if status > http.StatusAccepted && err == nil {
		fmt.Println(fmt.Sprintf("could not delete resource '%s/%s-%s'", client.GUID, kind, guid))
		return nil
	}
	if err != nil {
		printError(c, err)
		return err
	}

	fmt.Println(fmt.Sprintf("successfully delete resource '%s/%s-%s'", client.GUID, kind, guid))
	return nil
}
