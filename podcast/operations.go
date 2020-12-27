package podcast

import (
	"fmt"

	t "github.com/podops/podops/internal/types"
)

const (
	// AuthenticationRoute is used to verify a token
	authenticationRoute = "/_a/token"
	// productionRoute route to call ProductionEndpoint
	productionRoute = "/new"
	// resourceRoute route to call ResourceEndpoint
	resourceRoute = "/update/%s/%s/%s?f=%v" // "/update/:parent/:rsrc/:id"
	// listRoute route to call ListEndpoint
	listRoute = "/list"
	// buildRoute route to call BuildEndpoint
	buildRoute = "/build"
)

// CreateProduction invokes the CreateProductionEndpoint
func (cl *Client) CreateProduction(name, title, summary string) (*t.ProductionResponse, error) {

	if name == "" {
		return nil, fmt.Errorf("resource: name must not be empty")
	}

	req := t.ProductionRequest{
		Name:    name,
		Title:   title,
		Summary: summary,
	}

	resp := t.ProductionResponse{}
	_, err := cl.Post(cl.apiNamespace+productionRoute, &req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// List retrieves a list of resources
func (cl *Client) List() (*t.ProductionsResponse, error) {
	var resp t.ProductionsResponse

	_, err := cl.Get(cl.apiNamespace+listRoute, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateResource invokes the ResourceEndpoint
func (cl *Client) CreateResource(kind, guid string, force bool, rsrc interface{}) (int, error) {

	resp := t.StatusObject{}
	status, err := cl.Post(cl.apiNamespace+fmt.Sprintf(resourceRoute, cl.GUID, kind, guid, force), rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// UpdateResource invokes the ResourceEndpoint
func (cl *Client) UpdateResource(kind, guid string, force bool, rsrc interface{}) (int, error) {

	resp := t.StatusObject{}
	status, err := cl.Put(cl.apiNamespace+fmt.Sprintf(resourceRoute, cl.GUID, kind, guid, force), rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// Build invokes the BuildEndpoint
func (cl *Client) Build(guid string) (string, error) {
	req := t.BuildRequest{
		GUID: cl.GUID,
	}
	resp := t.BuildResponse{}

	_, err := cl.Post(cl.apiNamespace+buildRoute, &req, &resp)
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}
