package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/urfave/cli/v2"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/cli/netrc"
	"github.com/podops/podops/internal/loader"
)

const (
	machineEntry = "api.podops.dev"
)

var (
	client *podops.Client
)

func init() {
	opts := podops.LoadConfiguration()

	c, err := podops.NewClient(context.TODO(), opts.Token, opts)
	if err != nil {
		log.Fatal(err)
	}
	if c != nil {
		c.SetProduction(opts.Production)
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
	req.Header.Set("User-Agent", podops.UserAgentString)
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

func loadNetrc() *netrc.Netrc {
	nrc, _ := netrc.ParseFile(podops.DefaultConfigPath()) // FIXME test this, can we ignore err?
	if nrc == nil {
		nrc = &netrc.Netrc{}
	}
	return nrc
}

func storeLogin(userID, token string) error {
	nrc := loadNetrc()
	m := nrc.FindMachine(machineEntry)
	if m == nil {
		nrc.NewMachine(machineEntry, userID, token, "GUID")
	} else {
		m.UpdateLogin(userID)
		m.UpdatePassword(token)
	}
	data, _ := nrc.MarshalText()
	return ioutil.WriteFile(podops.DefaultConfigPath(), data, 0644)
}

func clearLogin() error {
	nrc := loadNetrc()
	m := nrc.FindMachine(machineEntry)
	if m != nil {
		nrc.RemoveMachine(machineEntry)

		data, _ := nrc.MarshalText()
		return ioutil.WriteFile(podops.DefaultConfigPath(), data, 0644)
	}
	return nil
}

func storeDefaultProduction(production string) error {
	nrc := loadNetrc()
	m := nrc.FindMachine(machineEntry)
	if m == nil {
		nrc.NewMachine(machineEntry, "", "", production)
	} else {
		m.UpdateAccount(production)
	}
	data, _ := nrc.MarshalText()
	return ioutil.WriteFile(podops.DefaultConfigPath(), data, 0644)
}

// GITHUB_ISSUE #15
func loadResource(path string) (interface{}, string, string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", "", err
	}

	r, kind, guid, err := loader.UnmarshalResource(data)
	if err != nil {
		return nil, "", "", err
	}

	return r, kind, guid, nil
}

func dumpResource(path string, doc interface{}) error {
	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, data, 0644)
	fmt.Printf("\n---\n# %s\n%s\n", path, string(data))

	return nil
}

func getProduction(c *cli.Context) string {
	prod := c.String("prod")
	if prod == "" {
		prod = client.DefaultProduction()
	}
	return prod
}

func productionListing(guid, name, title string, current bool) string {
	if current {
		return fmt.Sprintf("* %-20s%-50s%s", guid, name, title)
	}
	return fmt.Sprintf("  %-20s%-50s%s", guid, name, title)
}

func assetListing(guid, name, kind string) string {
	return fmt.Sprintf("  %-20s%-50s%s", guid, name, kind)
}

// printError formats a CLI error and prints it
func printError(c *cli.Context, err error) {
	msg := fmt.Sprintf("%s: %v", c.Command.Name, strings.ToLower(err.Error()))
	fmt.Println(msg)
}
