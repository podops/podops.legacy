package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/pkg/api"
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
		Namespace       string
	}
)

// DefaultClient returns the default configuration
func DefaultClient(token string) *Client {
	return &Client{
		ServiceEndpoint: a.DefaultAPIEndpoint,
		Token:           token,
		GUID:            "",
		Namespace:       api.NamespacePrefix,
	}
}

// Store persists the Client state
func (cl *Client) Store(path string) error {
	config, _ := json.Marshal(cl)

	// create the location if it does not exist
	baseDir := filepath.Dir(path)
	if baseDir != "." && baseDir != ".." {
		// path is contains a location
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(path, config, 0644)
}

// Validate verifies the token against the backend service
func (cl *Client) Validate() error {
	if cl.Token == "" {
		return a.ErrNoToken
	}
	status, err := cl.get(authenticationRoute, nil)
	if err != nil {
		return err
	}
	if status != http.StatusAccepted {
		// the only valid positive response
		return a.ErrNotAuthorized
	}
	return nil
}

// HasToken verifies that remote commands can be executed
func (cl *Client) HasToken() error {
	if cl.Token == "" {
		return a.ErrNotAuthorized
	}
	return nil
}

// HasTokenAndGUID verifies the presence of a token and GUID
func (cl *Client) HasTokenAndGUID() error {
	if cl.Token == "" {
		return a.ErrNotAuthorized
	}
	if cl.GUID == "" {
		return a.ErrInvalidParameters
	}
	return nil
}

// Get is used to request data from the API. No payload, only queries!
func (cl *Client) get(cmd string, response interface{}) (int, error) {
	url := cl.ServiceEndpoint + cmd

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return cl.invoke(req, response)
}

// Post is used to invoke an API method using http POST
func (cl *Client) post(cmd string, request, response interface{}) (int, error) {
	url := cl.ServiceEndpoint + cmd

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
	url := cl.ServiceEndpoint + cmd

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
	url := cl.ServiceEndpoint + cmd

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
	if cl.Token != "" {
		req.Header.Set("Authorization", "Bearer "+cl.Token)
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
	req.Header.Set("Authorization", "Bearer "+cl.Token)

	return req, err
}
