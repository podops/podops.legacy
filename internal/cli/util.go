package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/commons/pkg/errors"
)

var endpoint string = env.GetString("API_ENDPOINT", "https://api.podops.dev/a/v1")

// Post is used to invoke a CLI API method by posting a JSON payload.
func Post(cmd, token string, request, response interface{}) (int, error) {
	url := endpoint + cmd

	m, err := json.Marshal(&request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(m))
	if err != nil {
		return http.StatusBadRequest, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	// post the request to Slack
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, errors.New(fmt.Sprintf("API error: %s", resp.Status))
	}

	// unmarshal the response
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, err
}

// IsAuthorized does a quick verification
func IsAuthorized() bool {
	return DefaultValuesCLI.Token != ""
}
