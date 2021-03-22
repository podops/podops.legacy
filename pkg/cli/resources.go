package cli

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	"github.com/fupas/commons/pkg/util"

	a "github.com/podops/podops/apiv1"
)

// GetResourcesCommand list all resource associated with a show
func GetResourcesCommand(c *cli.Context) error {

	prod := getProduction(c)
	kind := resourceMap[strings.ToLower(c.Args().First())]

	if c.NArg() < 2 {
		// get a list of resources
		l, err := client.Resources(prod, kind)
		if err != nil {
			printError(c, err)
			return nil
		}

		if len(l.Resources) == 0 {
			fmt.Println("No resources found.")
		} else {
			fmt.Println(assetListing("ID", "NAME", "KIND"))
			for _, details := range l.Resources {
				fmt.Println(assetListing(details.GUID, details.Name, details.Kind))
			}
		}
	} else {
		// get a single resource
		guid := c.Args().Get(1)

		var rsrc interface{}
		err := client.GetResource(prod, kind, guid, &rsrc)
		if err != nil {
			printError(c, err)
			return nil
		}

		// FIXME verify that rsrc.Kind == kind

		data, err := yaml.Marshal(rsrc)
		if err != nil {
			return err
		}

		fmt.Printf("\n--- %s/%s-%s:\n\n%s\n\n", prod, kind, guid, string(data))
	}

	return nil
}

// CreateCommand creates a resource from a file, directory or URL
func CreateCommand(c *cli.Context) error {

	if c.NArg() != 1 {
		return fmt.Errorf("wrong number of arguments: expected 1, got %d", c.NArg())
	}
	path := c.Args().First()
	force := c.Bool("force")

	r, kind, guid, err := loadResource(path)
	if err != nil {
		return err
	}

	_, err = client.CreateResource(getProduction(c), kind, guid, force, r)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("created resource %s-%s", kind, guid))
	return nil
}

// UpdateCommand updates a resource from a file, directory or URL
func UpdateCommand(c *cli.Context) error {

	if c.NArg() != 1 {
		return fmt.Errorf("wrong number of arguments: expected 1, got %d", c.NArg())
	}
	path := c.Args().First()
	force := c.Bool("force")

	r, kind, guid, err := loadResource(path)
	if err != nil {
		return err
	}

	_, err = client.UpdateResource(getProduction(c), kind, guid, force, r)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Updated resource %s-%s", kind, guid))
	return nil
}

// DeleteResourcesCommand deletes a resource
func DeleteResourcesCommand(c *cli.Context) error {

	if c.NArg() != 2 {
		return fmt.Errorf("wrong number of arguments: expected 2, got %d", c.NArg())
	}

	prod := getProduction(c)
	kind := strings.ToLower(c.Args().First())
	guid := c.Args().Get(1)

	status, err := client.DeleteResource(prod, kind, guid)
	if err != nil {
		printError(c, err)
		return err
	}

	if status != http.StatusNoContent {
		fmt.Println(fmt.Sprintf("could not delete resource '%s/%s-%s'", prod, kind, guid))
		return nil
	}

	fmt.Println(fmt.Sprintf("successfully delete resource '%s/%s-%s'", prod, kind, guid))
	return nil
}

// TemplateCommand creates a resource template with all default values
func TemplateCommand(c *cli.Context) error {
	template := c.Args().First()
	if template != a.ResourceShow && template != a.ResourceEpisode {
		fmt.Println(fmt.Sprintf("\nDon't know how to create '%s'", template))
		return nil
	}

	parent := "PARENT-NAME"

	name := "NAME"
	if c.NArg() == 2 {
		name = c.Args().Get(1)
	}
	guid := c.String("guid")
	if guid == "" {
		guid, _ = util.ShortUUID()
	}
	parentGUID := c.String("parent")
	if parentGUID == "" {
		parentGUID = "PARENT-ID"
	}

	// create the yamls
	if template == "show" {
		show := a.DefaultShow(name, "TITLE", "SUMMARY", guid, a.DefaultPortalEndpoint, a.DefaultCDNEndpoint)
		err := dumpResource(fmt.Sprintf("show-%s.yaml", guid), show)
		if err != nil {
			printError(c, err)
			return nil
		}
	} else {
		episode := a.DefaultEpisode(name, parent, guid, parentGUID, a.DefaultPortalEndpoint, a.DefaultCDNEndpoint)
		err := dumpResource(fmt.Sprintf("episode-%s.yaml", guid), episode)
		if err != nil {
			printError(c, err)
			return nil
		}
	}

	return nil
}

// UploadCommand uploads an asset from a file
func UploadCommand(c *cli.Context) error {

	if c.NArg() != 1 {
		return fmt.Errorf("wrong number of arguments: expected 1, got %d", c.NArg())
	}

	prod := getProduction(c)
	name := c.Args().First()
	force := c.Bool("force")

	err := client.Upload(prod, name, force)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Uploaded '%s'", name))
	return nil
}
