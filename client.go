package podops

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	a "github.com/podops/podops/apiv1"
	cl "github.com/podops/podops/client"
	"github.com/podops/podops/pkg/api"
)

// NewClient creates a new podcast client.
//
// Clients should be reused instead of created as needed.
// The methods of a client instance are threadsafe.
func NewClient(token string) (*cl.Client, error) {
	client := cl.DefaultClient(token)
	if token != "" {
		if err := client.Validate(); err != nil {
			return nil, err
		}
	}
	return client, nil
}

// NewClientFromFile creates a client by reading values from a file
//
// Clients should be reused instead of created as needed. The methods of Client
// are safe for concurrent use by multiple goroutines.
func NewClientFromFile(path string) (*cl.Client, error) {
	var client *cl.Client

	if _, err := os.Stat(path); os.IsNotExist(err) {
		client = cl.DefaultClient("")
	} else {
		jsonFile, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &client)

		client.Namespace = api.NamespacePrefix
		client.ServiceEndpoint = a.DefaultAPIEndpoint
	}
	return client, nil
}

// DefaultConfigLocation returns the suggested default location for the config file
func DefaultConfigLocation() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".po/config")
}
