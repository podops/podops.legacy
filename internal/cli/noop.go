package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

// NoopCommand does nothing
func NoopCommand(c *cli.Context) error {
	fmt.Println(fmt.Sprintf("%s(%s) '%s %s' is not yet implemented!\n", CmdLineName, CmdLineVersion, CmdLineName, c.Command.Name))
	return nil
}
