package podops

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

const (
	// NamespacePrefix namespace for the client and CLI
	NamespacePrefix = "/a/v1"

	// productionRoute route to call ProductionEndpoint
	productionRoute = NamespacePrefix + "/production"
	// listProductionsRoute route to call ListProductionsEndpoint
	listProductionsRoute = NamespacePrefix + "/productions"

	// resourceRoute route to call ResourceEndpoint
	findResourceRoute   = NamespacePrefix + "/resource/%s"            // "/get/:id"
	getResourceRoute    = NamespacePrefix + "/resource/%s/%s/%s"      // "/get/:prod/:kind/:id"
	updateResourceRoute = NamespacePrefix + "/resource/%s/%s/%s?f=%v" // "/update/:prod/:kind/:id"
	listResourcesRoute  = NamespacePrefix + "/resource/%s/%s"
	deleteResourceRoute = NamespacePrefix + "/resource/%s/%s/%s"

	// buildRoute route to call BuildEndpoint
	buildRoute = NamespacePrefix + "/build"
	// uploadRoute route to UploadEndpoint
	uploadRoute = NamespacePrefix + "/upload"
)

// CreateProduction invokes the CreateProductionEndpoint
func (cl *Client) CreateProduction(name, title, summary string) (*Production, error) {
	if !cl.IsValid() {
		return nil, PodopsClientConfigurationErr
	}

	if name == "" {
		return nil, fmt.Errorf("name must not be empty") // FIXME replace with const
	}

	req := Production{
		Name:    name,
		Title:   title,
		Summary: summary,
	}

	resp := Production{}
	_, err := post(cl.opts.APIEndpoint, productionRoute, cl.opts.Token, &req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Productions retrieves a list of productions
func (cl *Client) Productions() (*ProductionList, error) {
	if !cl.IsValid() {
		return nil, PodopsClientConfigurationErr
	}

	var resp ProductionList
	_, err := get(cl.opts.APIEndpoint, listProductionsRoute, cl.opts.Token, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateResource invokes the ResourceEndpoint
func (cl *Client) CreateResource(production, kind, guid string, force bool, rsrc interface{}) (int, error) {
	if !cl.IsValid() {
		return http.StatusBadRequest, PodopsClientConfigurationErr
	}

	resp := StatusObject{}
	status, err := post(cl.opts.APIEndpoint, fmt.Sprintf(updateResourceRoute, production, kind, guid, force), cl.opts.Token, rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// GetResource returns a resource file
func (cl *Client) GetResource(production, kind, guid string, rsrc interface{}) error {
	if !cl.IsValid() {
		return PodopsClientConfigurationErr
	}

	status, err := get(cl.opts.APIEndpoint, fmt.Sprintf(getResourceRoute, production, kind, guid), cl.opts.Token, rsrc)
	if status == http.StatusBadRequest {
		return fmt.Errorf("not found: '%s/%s-%s'", production, kind, guid)
	}
	if err != nil {
		return err
	}

	return nil
}

// FindResource returns a resource file
func (cl *Client) FindResource(guid string, rsrc interface{}) error {
	if !cl.IsValid() {
		return PodopsClientConfigurationErr
	}

	status, err := get(cl.opts.APIEndpoint, fmt.Sprintf(findResourceRoute, guid), cl.opts.Token, rsrc)
	if status == http.StatusBadRequest {
		return fmt.Errorf("not found: '%s'", guid)
	}
	if err != nil {
		return err
	}

	return nil
}

// Resources retrieves a list of resources
func (cl *Client) Resources(production, kind string) (*ResourceList, error) {
	if !cl.IsValid() {
		return nil, PodopsClientConfigurationErr
	}
	if kind == "" {
		kind = "ALL"
	}

	var resp ResourceList
	_, err := get(cl.opts.APIEndpoint, fmt.Sprintf(listResourcesRoute, production, kind), cl.opts.Token, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateResource invokes the ResourceEndpoint
func (cl *Client) UpdateResource(production, kind, guid string, force bool, rsrc interface{}) (int, error) {
	if !cl.IsValid() {
		return http.StatusBadRequest, PodopsClientConfigurationErr
	}

	resp := StatusObject{}
	status, err := put(cl.opts.APIEndpoint, fmt.Sprintf(updateResourceRoute, production, kind, guid, force), cl.opts.Token, rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// DeleteResource deletes a resources
func (cl *Client) DeleteResource(production, kind, guid string) (int, error) {
	if !cl.IsValid() {
		return http.StatusBadRequest, PodopsClientConfigurationErr
	}

	status, err := delete(cl.opts.APIEndpoint, fmt.Sprintf(deleteResourceRoute, production, kind, guid), cl.opts.Token, nil)
	if err != nil {
		return status, err
	}
	return status, nil
}

// Build invokes the BuildEndpoint
func (cl *Client) Build(production string) (*BuildRequest, error) {
	if !cl.IsValid() {
		return nil, PodopsClientConfigurationErr
	}

	req := BuildRequest{
		GUID: production,
	}
	resp := BuildRequest{}

	_, err := post(cl.opts.APIEndpoint, buildRoute, cl.opts.Token, &req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Upload invokes the UploadEndpoint
func (cl *Client) Upload(production, path string, force bool) error {
	if !cl.IsValid() {
		return PodopsClientConfigurationErr
	}
	// FIXME this should be cl.opts.StorageEndpoint or CDNEndpoint
	req, err := upload(cl.opts.APIEndpoint, uploadRoute, cl.opts.Token, production, "asset", path)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", UserAgentString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	if resp.StatusCode > http.StatusNoContent {
		return fmt.Errorf("error uploading '%s': %s", path, resp.Status)
	}

	return nil
}
