package client

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	a "github.com/podops/podops/apiv1"
)

const (
	// AuthenticationRoute is used to verify a token
	authenticationRoute = "/_a/token"

	// productionRoute route to call ProductionEndpoint
	productionRoute = "/production"
	// listProductionsRoute route to call ListProductionsEndpoint
	listProductionsRoute = "/productions"

	// resourceRoute route to call ResourceEndpoint
	getResourceRoute    = "/resource/%s/%s/%s"      // "/update/:prod/:kind/:id"
	updateResourceRoute = "/resource/%s/%s/%s?f=%v" // "/update/:prod/:kind/:id"
	listResourcesRoute  = "/resource/%s/%s"
	deleteResourceRoute = "/resource/%s/%s/%s"

	// buildRoute route to call BuildEndpoint
	buildRoute = "/build"
	// uploadRoute route to UploadEndpoint
	uploadRoute = "/upload"
)

// CreateProduction invokes the CreateProductionEndpoint
func (cl *Client) CreateProduction(name, title, summary string) (*a.Production, error) {
	if !cl.Valid() {
		return nil, PodopsClientConfigurationErr
	}

	if name == "" {
		return nil, fmt.Errorf("name must not be empty") // FIXME replace with const
	}

	req := a.Production{
		Name:    name,
		Title:   title,
		Summary: summary,
	}

	resp := a.Production{}
	_, err := cl.post(cl.ns+productionRoute, &req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Productions retrieves a list of productions
func (cl *Client) Productions() (*a.ProductionList, error) {
	if !cl.Valid() {
		return nil, PodopsClientConfigurationErr
	}

	var resp a.ProductionList
	_, err := cl.get(cl.ns+listProductionsRoute, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateResource invokes the ResourceEndpoint
func (cl *Client) CreateResource(prod, kind, rsrcGUID string, force bool, rsrc interface{}) (int, error) {
	if !cl.Valid() {
		return http.StatusBadRequest, PodopsClientConfigurationErr
	}

	resp := a.StatusObject{}
	status, err := cl.post(cl.ns+fmt.Sprintf(updateResourceRoute, prod, kind, rsrcGUID, force), rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// GetResource returns a resource file
func (cl *Client) GetResource(prod, kind, guid string, rsrc interface{}) error {
	if !cl.Valid() {
		return PodopsClientConfigurationErr
	}

	status, err := cl.get(cl.ns+fmt.Sprintf(getResourceRoute, prod, kind, guid), rsrc)
	if status == http.StatusBadRequest {
		return fmt.Errorf("not found: '%s/%s-%s'", prod, kind, guid)
	}
	if err != nil {
		return err
	}

	return nil
}

// UpdateResource invokes the ResourceEndpoint
func (cl *Client) UpdateResource(prod, kind, rsrcGUID string, force bool, rsrc interface{}) (int, error) {
	if !cl.Valid() {
		return http.StatusBadRequest, PodopsClientConfigurationErr
	}

	resp := a.StatusObject{}
	status, err := cl.put(cl.ns+fmt.Sprintf(updateResourceRoute, prod, kind, rsrcGUID, force), rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// Resources retrieves a list of resources
func (cl *Client) Resources(prod, kind string) (*a.ResourceList, error) {
	if !cl.Valid() {
		return nil, PodopsClientConfigurationErr
	}
	if kind == "" {
		kind = "ALL"
	}

	var resp a.ResourceList
	_, err := cl.get(cl.ns+fmt.Sprintf(listResourcesRoute, prod, kind), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteResource deletes a resources
func (cl *Client) DeleteResource(prod, kind, guid string) (int, error) {
	if !cl.Valid() {
		return http.StatusBadRequest, PodopsClientConfigurationErr
	}

	status, err := cl.delete(cl.ns+fmt.Sprintf(deleteResourceRoute, prod, kind, guid), nil)
	if err != nil {
		return status, err
	}
	return status, nil
}

// Build invokes the BuildEndpoint
func (cl *Client) Build(prod string) (*a.Build, error) {
	if !cl.Valid() {
		return nil, PodopsClientConfigurationErr
	}

	req := a.Build{
		GUID: prod,
	}
	resp := a.Build{}

	_, err := cl.post(cl.ns+buildRoute, &req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Upload invokes the UploadEndpoint
func (cl *Client) Upload(prod, path string, force bool) error {
	if !cl.Valid() {
		return PodopsClientConfigurationErr
	}

	req, err := cl.fileUploadRequest(cl.opts.APIEndpoint+cl.ns+uploadRoute, prod, path) // FIXME upload should go against the CDN endpoint
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", a.UserAgentString)

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
