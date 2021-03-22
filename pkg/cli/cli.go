package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"

	a "github.com/podops/podops/apiv1"
)

const (
	// BasicCmdGroup groups basic commands
	BasicCmdGroup = "\nBasic Commands"
	// SettingsCmdGroup groups settings
	SettingsCmdGroup = "\nSettings Commands"
	// ShowCmdGroup groups basic show commands
	ShowCmdGroup = "\nContent Creation Commands"
	// ShowMgmtCmdGroup groups advanced show commands
	ShowMgmtCmdGroup = "\nContent Management Commands"
)

// NoOpCommand is just a placeholder
func NoOpCommand(c *cli.Context) error {
	return cli.Exit(fmt.Sprintf("Command '%s' is not implemented", c.Command.Name), 0)
}

// Post is used to invoke an API method using http POST
func post(url string, request, response interface{}) (int, error) {
	m, err := json.Marshal(&request)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(m))
	if err != nil {
		return http.StatusBadRequest, err
	}

	return invoke(req, response)
}

func invoke(req *http.Request, response interface{}) (int, error) {

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", a.UserAgentString)

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

	/*
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
	*/
	// unmarshal the response if one is expected
	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return resp.StatusCode, nil
}
