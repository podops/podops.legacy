package commands

import (
	"fmt"
)

func shorten(s string, l int) string {
	if len(s) <= l {
		return s
	}
	return fmt.Sprintf("%s..%s", s[0:(l-11)], s[(len(s)-9):])
}

func productionListing(guid, name, title string, current bool) string {
	if current {
		return fmt.Sprintf("* %-20s%-50s%s", guid, name, title)
	}
	return fmt.Sprintf("  %-20s%-50s%s", guid, name, title)
}

func assetListing(guid, name, kind string) string {
	return fmt.Sprintf("  %-20s%-50s%s", guid, name, kind)
}
