package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/pkg/api"
)

var (
	PodopsClientConfigurationErr error = fmt.Errorf("new client: invalid configuration")
)

// Client is a client for interacting with the PodOps service.
//
// Clients should be reused instead of created as needed.
// The methods of Client are safe for concurrent use by multiple goroutines.
type (
	ClientOption struct {
		Token          string
		Production     string
		APIEndpoint    string
		CDNEndpoint    string
		PortalEndpoint string
	}

	Client struct {
		opts              *ClientOption
		defaultProduction string
		// internal for now
		validated bool
		valid     bool
		ns        string
		realm     string
	}
)

func New(ctx context.Context, o *ClientOption) (*Client, error) {
	if o == nil || !o.Valid() {
		return nil, PodopsClientConfigurationErr
	}
	return &Client{
		opts:      o,
		validated: false,
		valid:     false,
		ns:        api.NamespacePrefix,
		realm:     "podops",
	}, nil
}

// Valid checks if all configuration parameters are provided
func (cl *Client) Valid() bool {
	if cl.validated {
		return cl.valid
	}
	cl.validated = true
	// verify the opts first
	if !cl.opts.Valid() {
		cl.valid = false
		return false
	}

	cl.valid = true // FIXME try to verify the token against the API

	return true
}

func (cl *Client) SetProduction(production string) {
	cl.defaultProduction = production
}

func (cl *Client) APIEndpoint() string {
	return cl.opts.APIEndpoint
}

func (cl *Client) CDNEndpoint() string {
	return cl.opts.CDNEndpoint
}

func (cl *Client) PortalEndpoint() string {
	return cl.opts.PortalEndpoint
}

func (cl *Client) Realm() string {
	return cl.realm
}

func (cl *Client) Token() string {
	return cl.opts.Token
}

func (cl *Client) DefaultProduction() string {
	return cl.defaultProduction
}

// Merge clones co and combines it with the provided options
func (co ClientOption) Merge(opts *ClientOption) *ClientOption {
	o := &ClientOption{}
	o.Token = co.Token
	o.APIEndpoint = co.APIEndpoint
	o.CDNEndpoint = co.CDNEndpoint
	o.PortalEndpoint = co.PortalEndpoint

	if opts != nil {
		if opts.Token != "" {
			o.Token = opts.Token
		}
		if opts.APIEndpoint != "" {
			o.APIEndpoint = opts.APIEndpoint
		}
		if opts.CDNEndpoint != "" {
			o.CDNEndpoint = opts.CDNEndpoint
		}
		if opts.PortalEndpoint != "" {
			o.PortalEndpoint = opts.PortalEndpoint
		}
	}

	return o
}

// Valid checks if all configuration parameters are provided
func (co ClientOption) Valid() bool {
	return co.APIEndpoint != "" && co.CDNEndpoint != "" && co.PortalEndpoint != ""
	// we can not validate token and production. There are some API calls that do not need them...
}

// Get is used to request data from the API. No payload, only queries!
func (cl *Client) get(cmd string, response interface{}) (int, error) {
	url := cl.opts.APIEndpoint + cmd

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return cl.invoke(req, response)
}

// Post is used to invoke an API method using http POST
func (cl *Client) post(cmd string, request, response interface{}) (int, error) {
	url := cl.opts.APIEndpoint + cmd

	m, err := json.Marshal(&request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(m))
	if err != nil {
		return http.StatusBadRequest, err
	}

	return cl.invoke(req, response)
}

// Put is used to invoke an API method using http PUT
func (cl *Client) put(cmd string, request, response interface{}) (int, error) {
	url := cl.opts.APIEndpoint + cmd

	m, err := json.Marshal(&request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(m))
	if err != nil {
		return http.StatusBadRequest, err
	}

	return cl.invoke(req, response)
}

// DELETE is used to request the deletion of a resource. Maybe apayload, no response!
func (cl *Client) delete(cmd string, request interface{}) (int, error) {
	url := cl.opts.APIEndpoint + cmd

	var req *http.Request
	var err error

	if request != nil {
		m, err := json.Marshal(&request)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		req, err = http.NewRequest("DELETE", url, bytes.NewBuffer(m))
	} else {
		req, err = http.NewRequest("DELETE", url, nil)
	}
	if err != nil {
		return http.StatusBadRequest, err
	}

	return cl.invoke(req, nil)
}

func (cl *Client) invoke(req *http.Request, response interface{}) (int, error) {

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", a.UserAgentString)
	if cl.opts.Token != "" {
		req.Header.Set("Authorization", "Bearer "+cl.opts.Token)
	}

	// perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if resp == nil {
			return http.StatusInternalServerError, err
		}
		return resp.StatusCode, err
	}

	defer resp.Body.Close()

	// anything other than OK, Created, Accepted, NoContent is treated as an error
	if resp.StatusCode > http.StatusNoContent {
		if response != nil {
			// as we expect a response, there might be a StatusObject
			status := &a.StatusObject{}
			err = json.NewDecoder(resp.Body).Decode(&status)
			if err != nil {
				return resp.StatusCode, fmt.Errorf("status: %d", resp.StatusCode)
			}
			return status.Status, fmt.Errorf(status.Message)
		}
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

// FIXME this implementation does not work for VERY large files !

// Creates a new file upload http request with optional extra params
func (cl *Client) fileUploadRequest(uri, guid, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("asset", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri+"/"+guid, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+cl.opts.Token)

	return req, err
}
