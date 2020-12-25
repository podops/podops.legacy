package cli

import (
	"os"
)

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
