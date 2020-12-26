package podcast

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/txsvc/commons/pkg/env"
)

const (
	// presetNameAndPath is the name and location of the config file
	presetNameAndPath = ".po"

	// DefaultServiceEndpoint is the service URL
	DefaultServiceEndpoint = "https://api.podops.dev"
	// DefaultNamespacePrefix is the API's namespace
	DefaultNamespacePrefix = "/a/v1"
)

// Client is a client for interacting with the PodOps service.
//
// Clients should be reused instead of created as needed.
// The methods of Client are safe for concurrent use by multiple goroutines.
type (
	Client struct {
		ServiceEndpoint string `json:"url" binding:"required"`
		Token           string `json:"token" binding:"required"`
		GUID            string `json:"guid" binding:"required"`
		authorized      bool
		apiNamespace    string
	}
)

// NewClient creates a new podcast client.
//
// Clients should be reused instead of created as needed. The methods of Client
// are safe for concurrent use by multiple goroutines.
func NewClient(ctx context.Context, token string) (*Client, error) {
	client := &Client{
		ServiceEndpoint: env.GetString("API_ENDPOINT", DefaultServiceEndpoint),
		Token:           token,
		GUID:            "",
		authorized:      false,
		apiNamespace:    DefaultNamespacePrefix,
	}
	if err := client.Validate(); err != nil {
		return nil, err
	}
	return client, nil
}

// NewClientFromFile creates a client by reading values from a file
//
// Clients should be reused instead of created as needed. The methods of Client
// are safe for concurrent use by multiple goroutines.
func NewClientFromFile(ctx context.Context, path string) (*Client, error) {
	var client *Client

	if _, err := os.Stat(path); os.IsNotExist(err) {
		client = &Client{
			ServiceEndpoint: env.GetString("API_ENDPOINT", DefaultServiceEndpoint),
			Token:           "",
			GUID:            "",
			authorized:      false,
			apiNamespace:    DefaultNamespacePrefix,
		}
		return client, nil
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &client)
	client.apiNamespace = DefaultNamespacePrefix

	if err := client.Validate(); err != nil {
		return nil, err
	}

	return client, nil
}

// Close does whatever kind of clean-up is necessary
func (cl *Client) Close() {
	// FIXME: just a placeholder for now
}

// Store persists the Client state
func (cl *Client) Store(path string) {
	defaults, _ := json.Marshal(cl)
	ioutil.WriteFile(path, defaults, 0644)
}

// IsAuthorized does a quick verification
func (cl *Client) IsAuthorized() bool {
	return cl.authorized
}

// Valid verifies that remote commands can be executed
func (cl *Client) Valid() error {
	if !cl.authorized {
		return fmt.Errorf("Not authorized. Use 'po auth' first") // FIXME generic text, not CLI specific
	}
	if cl.GUID == "" {
		return fmt.Errorf("No show selected. Use 'po show' first")
	}
	return nil
}

// Validate verifies the token against the backend service
func (cl *Client) Validate() error {

	if cl.Token == "" {
		return fmt.Errorf("validation: missing token")
	}

	status, err := cl.Get(authenticationRoute, nil)
	if err != nil {
		return err
	}
	if status != http.StatusAccepted {
		// the only valid positive response
		return fmt.Errorf("validation: not authorized %d", status)
	}
	cl.authorized = true
	return nil
}
