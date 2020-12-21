package cli

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	// presetNameAndPath is the name and location of the config file
	presetNameAndPath = ".po"

	// DefaultServiceEndpoint is the service URL
	DefaultServiceEndpoint = "https://api.podops.dev"
)

type (
	// DefaultValues stores all presets the CLI needs
	DefaultValues struct {
		ServiceEndpoint string `json:"url"`
		Token           string `json:"token"`
		ClientID        string `json:"client_id"`
		DefaultShow     string `json:"show"`
	}
)

var DefaultValuesCLI *DefaultValues

func init() {
	df := &DefaultValues{
		ServiceEndpoint: DefaultServiceEndpoint,
		Token:           "",
		ClientID:        "",
		DefaultShow:     "",
	}
	DefaultValuesCLI = df
}

// ServiceEndpoint returns the service endpoint
func ServiceEndpoint() string {
	return DefaultValuesCLI.ServiceEndpoint
}

// Token returns the API token of the current user
func Token() string {
	return DefaultValuesCLI.Token
}

// ClientID returns the users ID
func ClientID() string {
	return DefaultValuesCLI.ClientID
}

// DefaultShow returns the current show
func DefaultShow() string {
	return DefaultValuesCLI.DefaultShow
}

// LoadOrCreateConfig initializes the default settings
func LoadOrCreateConfig() {
	if _, err := os.Stat(presetNameAndPath); os.IsNotExist(err) {
		StoreConfig()
	} else {
		jsonFile, err := os.Open(presetNameAndPath)
		if err != nil {
			return
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, DefaultValuesCLI)
	}
}

// StoreConfig persists the current values
func StoreConfig() {
	defaults, _ := json.Marshal(DefaultValuesCLI)
	ioutil.WriteFile(presetNameAndPath, defaults, 0644)
}
