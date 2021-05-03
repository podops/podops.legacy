package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/txsvc/platform/v2/pkg/api"

	"github.com/podops/podops/internal/messagedef"
)

/*
Make sure to update version numbers in these locations also:

- version.go
- .github/*
*/

const (
	// MajorVersion of the API
	majorVersion = 1
	// MinorVersion of the API
	minorVersion = 0
	// FixVersion of the API
	fixVersion = 2
)

var (
	// UserAgentString identifies any http request podops makes
	UserAgentString string = fmt.Sprintf("PodOps %d.%d.%d", majorVersion, minorVersion, fixVersion)
)

// Get is used to request data from the API. No payload, only queries!
func Get(url, cmd, token string, response interface{}) (int, error) {
	uri := url + cmd

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return invoke(token, req, response)
}

// Post is used to invoke an API method using http POST
func Post(url, cmd, token string, request, response interface{}) (int, error) {
	uri := url + cmd

	m, err := json.Marshal(&request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(m))
	if err != nil {
		return http.StatusBadRequest, err
	}

	return invoke(token, req, response)
}

// Put is used to invoke an API method using http PUT
func Put(url, cmd, token string, request, response interface{}) (int, error) {
	uri := url + cmd

	m, err := json.Marshal(&request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(m))
	if err != nil {
		return http.StatusBadRequest, err
	}

	return invoke(token, req, response)
}

// DELETE is used to request the deletion of a resource. Maybe apayload, no response!
func Delete(url, cmd, token string, request interface{}) (int, error) {
	uri := url + cmd

	if request != nil {
		m, err := json.Marshal(&request)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		req, err := http.NewRequest("DELETE", uri, bytes.NewBuffer(m))
		if err != nil {
			return http.StatusBadRequest, err
		}
		return invoke(token, req, nil)
	}

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return invoke(token, req, nil)

}

// GITHUB_ISSUE #13

// Creates a new file upload http request with optional extra params
func Upload(url, cmd, token, guid, form, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile(form, filepath.Base(path))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	uri := url + cmd + "/" + guid
	req, err := http.NewRequest("POST", uri, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", UserAgentString)

	return req, err
}

func invoke(token string, req *http.Request, response interface{}) (int, error) {

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", UserAgentString)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
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
			status := api.StatusObject{}
			err = json.NewDecoder(resp.Body).Decode(&status)
			if err != nil {
				return resp.StatusCode, fmt.Errorf(messagedef.MsgStatus, resp.StatusCode)
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
