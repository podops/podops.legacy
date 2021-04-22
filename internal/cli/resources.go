package cli

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	"github.com/fupas/commons/pkg/util"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/metadata"
)

// GetResourcesCommand list all resource associated with a show
func GetResourcesCommand(c *cli.Context) error {
	kind := podops.ResourceALL
	prod := getProduction(c)

	if c.NArg() == 1 {
		k, err := backend.NormalizeKind(c.Args().First())
		if err == nil {
			kind = k
		} else {
			kind = ""
		}
	}

	if kind != "" {
		// get a list of resources
		if kind == "" {
			kind = podops.ResourceALL
		}

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
				if details.Kind == podops.ResourceAsset {
					name := "???"

					if details.EnclosureURI != "" {
						name = metadata.LocalNamePart(details.EnclosureURI)
					} else if details.ImageURI != "" {
						name = metadata.LocalNamePart(details.ImageURI)
					}

					fmt.Println(assetListing(details.GUID, name, details.Kind))
				} else {
					fmt.Println(assetListing(details.GUID, details.Name, details.Kind))
				}
			}
		}
	} else {
		guid := c.Args().First()

		var rsrc interface{}
		err := client.FindResource(guid, &rsrc)
		if err != nil {
			fmt.Println("Resource not found.")
			return nil
		}

		data, err := yaml.Marshal(rsrc)
		if err != nil {
			return err
		}

		fmt.Printf("\n---\n# %s/%s\n%s\n\n", prod, guid, string(data))
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

	fmt.Println(fmt.Sprintf("successfully deleted resource '%s/%s-%s'", prod, kind, guid))
	return nil
}

// TemplateCommand creates a resource template with all default values
func TemplateCommand(c *cli.Context) error {
	template := c.Args().First()
	if template != podops.ResourceShow && template != podops.ResourceEpisode {
		fmt.Println(fmt.Sprintf("\nDon't know how to create '%s'", template))
		return nil
	}

	parentName := "PARENT-NAME"

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
		show := podops.DefaultShow(name, "TITLE", "SUMMARY", guid, podops.DefaultEndpoint, podops.DefaultCDNEndpoint)
		err := dumpResource(fmt.Sprintf("show-%s.yaml", guid), show)
		if err != nil {
			printError(c, err)
			return nil
		}
	} else {
		episode := podops.DefaultEpisode(name, parentName, guid, parentGUID, podops.DefaultEndpoint, podops.DefaultCDNEndpoint)
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
