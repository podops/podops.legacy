package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"path/filepath"

	"github.com/fupas/commons/pkg/env"
	"github.com/urfave/cli/v2"

	"github.com/podops/podops"
	a "github.com/podops/podops/apiv1"
	cl "github.com/podops/podops/pkg/client"
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

	machineEntry = "api.podops.dev"
)

var (
	client *cl.Client
)

func init() {
	cl := podops.DefaultClientOptions()

	nrc := loadNetrc()
	m := nrc.FindMachine(machineEntry)
	if m != nil {
		cl.Token = m.Password
	}

	c, err := podops.NewClient(context.TODO(), cl.Token, cl)
	if err != nil {
		log.Fatal(err)
	}
	if c != nil {
		client = c
	}
}

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
	if client.Token() != "" {
		req.Header.Set("Authorization", "Bearer "+client.Token())
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

	// unmarshal the response if one is expected
	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return resp.StatusCode, nil
}

func netrcPath() string {
	path := env.GetString("PODOPS_CREDENTIALS", "")
	if path == "" {
		usr, _ := user.Current()
		path = filepath.Join(usr.HomeDir, ".netrc")
	}
	return path
}

func loadNetrc() *Netrc {
	nrc, _ := ParseFile(netrcPath()) // FIXME test this, can we ignore err?
	if nrc == nil {
		nrc = &Netrc{machines: make([]*Machine, 0, 20), macros: make(Macros, 10)}
	}
	return nrc
}

func updateNetrc(userID, clientID, token string) error {
	nrc := loadNetrc()
	m := nrc.FindMachine(machineEntry)
	if m == nil {
		m = nrc.NewMachine(machineEntry, userID, token, clientID)
	} else {
		m.UpdateLogin(userID)
		m.UpdatePassword(token)
	}
	data, _ := nrc.MarshalText()
	return ioutil.WriteFile(netrcPath(), data, 0644)
}
