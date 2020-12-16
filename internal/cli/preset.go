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

var defaultValues *DefaultValues

func init() {
	df := &DefaultValues{
		ServiceEndpoint: DefaultServiceEndpoint,
		Token:           "",
		ClientID:        "",
		DefaultShow:     "",
	}
	defaultValues = df
}

// ServiceEndpoint returns the service endpoint
func ServiceEndpoint() string {
	return defaultValues.ServiceEndpoint
}

// Token returns the API token of the current user
func Token() string {
	return defaultValues.Token
}

// ClientID returns the users ID
func ClientID() string {
	return defaultValues.ClientID
}

// DefaultShow returns the current show
func DefaultShow() string {
	return defaultValues.DefaultShow
}

// LoadOrCreateDefaultValues initializes the default settings
func LoadOrCreateDefaultValues() {
	if _, err := os.Stat(presetNameAndPath); os.IsNotExist(err) {
		StoreDefaultValues()
	} else {
		jsonFile, err := os.Open(presetNameAndPath)
		if err != nil {
			return
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, defaultValues)
	}
}

// StoreDefaultValues persists the current values
func StoreDefaultValues() {
	defaults, _ := json.Marshal(defaultValues)
	ioutil.WriteFile(presetNameAndPath, defaults, 0644)
}
