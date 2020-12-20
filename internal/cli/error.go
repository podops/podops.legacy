package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"

	"github.com/urfave/cli"
)

func PrintError(c *cli.Context, err error) {
	fmt.Println(fmt.Errorf("%w", err))
}
