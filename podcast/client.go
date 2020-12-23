package podcast

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/txsvc/commons/pkg/env"
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

	if err := client.Validate(); err != nil {
		return nil, err
	}

	return client, nil
}

// Close does whatever kind of clean-up is necessary
func (cl *Client) Close() {
	// FIXME: just a placeholder
}
