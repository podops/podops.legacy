package podcast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/podops/podops/internal/errors"
)

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
