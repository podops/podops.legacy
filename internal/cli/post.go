package cli

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/txsvc/commons/pkg/env"
)

var endpoint string = env.GetString("API_ENDPOINT", "https://api.podops.dev/a/1/cli")

// Post is used to invoke a CLI API method by posting a JSON payload.
func Post(cmd, token string, request, response interface{}) error {
	url := endpoint + cmd

	m, err := json.Marshal(&request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(m))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	// post the request to Slack
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// unmarshal the response
	err = json.NewDecoder(resp.Body).Decode(response)

	return err
}
