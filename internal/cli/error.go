package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"

	"github.com/urfave/cli"
)

// PrintError formats a CLI error and prints it
func PrintError(c *cli.Context, err error) {
	fmt.Println(fmt.Errorf("%w", err))
}
