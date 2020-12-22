package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/txsvc/commons/pkg/env"
	"github.com/urfave/cli"

	"github.com/podops/podops/internal/errors"
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

	return resp.StatusCode, err
}

// IsAuthorized does a quick verification
func IsAuthorized() bool {
	return DefaultValuesCLI.Token != ""
}

// PrintError formats a CLI error and prints it
func PrintError(c *cli.Context, operation string, status int, err error) {
	msg := ""
	switch status {
	case http.StatusInternalServerError:
		msg = fmt.Sprintf("Oops, something went wrong! [%s]", operation)
		break
	case http.StatusConflict:
		msg = fmt.Sprintf("Could not create resource [%s]", operation)
		break
	default:
		msg = err.Error()
	}
	fmt.Println(msg)
}
