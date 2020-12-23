package podcast

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/podops/podops/internal/errors"
	t "github.com/podops/podops/pkg/types"
)

const (
	// presetNameAndPath is the name and location of the config file
	presetNameAndPath = ".po"

	// DefaultServiceEndpoint is the service URL
	DefaultServiceEndpoint = "https://api.podops.dev"

	// AuthenticationRoute is used to verify a token
	authenticationRoute = "/_a/token"
	// productionRoute route to call ProductionEndpoint
	productionRoute = "/a/v1/new"
	// resourceRoute route to call ResourceEndpoint
	resourceRoute = "/a/v1/update/%s/%s/%s" // "/update/:parent/:rsrc/:id"
)

// CreateProduction invokes the CreateProductionEndpoint
func (cl *Client) CreateProduction(name, title, summary string) (*t.ProductionResponse, error) {

	if name == "" {
		return nil, errors.New("Name must not be empty", http.StatusBadRequest)
	}

	req := t.ProductionRequest{
		Name:    name,
		Title:   title,
		Summary: summary,
	}

	resp := t.ProductionResponse{}
	_, err := cl.Post(productionRoute, &req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateResource invokes the ResourceEndpoint
func (cl *Client) UpdateResource(kind, guid string, rsrc interface{}) (int, error) {

	resp := errors.StatusObject{}
	status, err := cl.Post(fmt.Sprintf(resourceRoute, cl.GUID, kind, guid), rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// Store persists the Client
func (cl *Client) Store(path string) {
	defaults, _ := json.Marshal(cl)
	ioutil.WriteFile(path, defaults, 0644)
}

// IsAuthorized does a quick verification
func (cl *Client) IsAuthorized() bool {
	return cl.Token != ""
}

// Validate verifies the token against the backend service
func (cl *Client) Validate() error {

	if cl.Token == "" {
		return errors.New("Missing token", http.StatusBadRequest)
	}

	status, err := cl.Get(authenticationRoute, nil)
	if err != nil {
		return err
	}
	if status != http.StatusAccepted {
		// the only valid positive response
		return errors.New("Not authorized", status)
	}
	cl.authorized = true
	return nil
}
