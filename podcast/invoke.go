package podcast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	t "github.com/podops/podops/internal/types"
)

// Get is used to request data from the API. No payload, only queries!
func (cl *Client) Get(cmd string, response interface{}) (int, error) {
	url := cl.ServiceEndpoint + cmd

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return cl.invoke(req, response)
}

// Post is used to invoke an API method using http POST
func (cl *Client) Post(cmd string, request, response interface{}) (int, error) {
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
func (cl *Client) Put(cmd string, request, response interface{}) (int, error) {
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

func (cl *Client) invoke(req *http.Request, response interface{}) (int, error) {

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+cl.Token)

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
			status := &t.StatusObject{}
			err = json.NewDecoder(resp.Body).Decode(&status)
			if err != nil {
				return resp.StatusCode, fmt.Errorf(fmt.Sprintf("status: %d", resp.StatusCode))
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

// FIXME does not work for VERY large files !

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
