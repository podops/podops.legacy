package podcast

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/podops/podops/internal/errors"
	"github.com/txsvc/commons/pkg/env"
)

const (
	// presetNameAndPath is the name and location of the config file
	presetNameAndPath = ".po"

	// DefaultServiceEndpoint is the service URL
	DefaultServiceEndpoint = "https://api.podops.dev"

	// AuthenticationRoute is used to verify a token
	authenticationRoute = "/_a/token"
	// productionRoute route to call ProductionEndpoint
	productionRoute = "/a/v1/new"
	// resourceRoute route to call ResourceEndpoint
	resourceRoute = "/a/v1/update/%s/%s" // "/update/:rsrc/:id"
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

// Store persists the Client
func (cl *Client) Store(path string) {
	defaults, _ := json.Marshal(cl)
	ioutil.WriteFile(path, defaults, 0644)
}

// IsAuthorized does a quick verification
func (cl *Client) IsAuthorized() bool {
	return cl.Token != ""
}

// Validate verifies the token against the backend service
func (cl *Client) Validate() error {

	if cl.Token == "" {
		return errors.New("Missing token", http.StatusBadRequest)
	}

	status, err := cl.Get(authenticationRoute, nil)
	if err != nil {
		return err
	}
	if status != http.StatusAccepted {
		// the only valid positive response
		return errors.New("Not authorized", status)
	}
	cl.authorized = true
	return nil
}

// Post is used to invoke an API method by posting a JSON payload.
func (cl *Client) Post(cmd string, request, response interface{}) (int, error) {

	m, err := json.Marshal(&request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("POST", cl.ServiceEndpoint+cmd, bytes.NewBuffer(m))
	if err != nil {
		return http.StatusBadRequest, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+cl.Token)

	// post the request to Slack
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp.StatusCode, err
	}

	defer resp.Body.Close()

	// anything other than OK, Created, Accepted, No Content is treated as an error
	if resp.StatusCode > http.StatusNoContent {
		return resp.StatusCode, errors.New(fmt.Sprintf("Status %d", resp.StatusCode), resp.StatusCode)
	}

	// FIXME: support empty body e.g. for StatusAccepted ...

	// unmarshal the response
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, nil
}

// Get is used to invoke an API method by requesting an URI. No payload, only queries!
func (cl *Client) Get(cmd string, response interface{}) (int, error) {

	req, err := http.NewRequest("GET", cl.ServiceEndpoint+cmd, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+cl.Token)

	// post the request to Slack
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp.StatusCode, err
	}

	defer resp.Body.Close()

	// anything other than OK, Created, Accepted, No Content is treated as an error
	if resp.StatusCode > http.StatusNoContent {
		return resp.StatusCode, errors.New(fmt.Sprintf("Status %d", resp.StatusCode), resp.StatusCode)
	}

	// unmarshal the response if one is expected
	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return resp.StatusCode, nil
}
