package podops

import (
	"context"

	"github.com/fupas/commons/pkg/env"

	a "github.com/podops/podops/apiv1"
	cl "github.com/podops/podops/pkg/client"
)

// NewClient creates a new podcast client.
//
// Clients should be reused instead of created as needed.
// The methods of a client instance are threadsafe.
func NewClient(ctx context.Context, token string, opts ...*cl.ClientOption) (*cl.Client, error) {

	co := DefaultClientOptions()
	if len(opts) != 0 {
		// FIXME we assume only 1 opts is provided!
		co = co.Merge(opts[0])
	}

	if token != "" {
		co.Token = token
	}

	return cl.New(ctx, co)
}

/*
// NewClientFromFile creates a client by reading values from a file
//
// Clients should be reused instead of created as needed. The methods of Client
// are safe for concurrent use by multiple goroutines.
func _NewClientFromFile(path string) (*cl.Client, error) {
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
*/

// DefaultClientOptions returns a default client configuration bases on ENV variables
func DefaultClientOptions() *cl.ClientOption {
	o := &cl.ClientOption{}
	o.Token = env.GetString("PODOPS_API_TOKEN", "")
	o.APIEndpoint = a.DefaultAPIEndpoint
	o.CDNEndpoint = a.DefaultCDNEndpoint
	return o
}
