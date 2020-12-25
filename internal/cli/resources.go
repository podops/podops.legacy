package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"

	"github.com/urfave/cli"
)

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
