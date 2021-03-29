package client

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/pkg/backend/models"
)

const (
	// AuthenticationRoute is used to verify a token
	authenticationRoute = "/_a/token"

	// productionRoute route to call ProductionEndpoint
	productionRoute = "/production"
	// listProductionsRoute route to call ListProductionsEndpoint
	listProductionsRoute = "/productions"

	// resourceRoute route to call ResourceEndpoint
	findResourceRoute   = "/resource/%s"            // "/get/:id"
	getResourceRoute    = "/resource/%s/%s/%s"      // "/get/:prod/:kind/:id"
	updateResourceRoute = "/resource/%s/%s/%s?f=%v" // "/update/:prod/:kind/:id"
	listResourcesRoute  = "/resource/%s/%s"
	deleteResourceRoute = "/resource/%s/%s/%s"

	// buildRoute route to call BuildEndpoint
	buildRoute = "/build"
	// uploadRoute route to UploadEndpoint
	uploadRoute = "/upload"
)

// CreateProduction invokes the CreateProductionEndpoint
func (cl *Client) CreateProduction(name, title, summary string) (*models.Production, error) {
	if !cl.IsValid() {
		return nil, PodopsClientConfigurationErr
	}

	if name == "" {
		return nil, fmt.Errorf("name must not be empty") // FIXME replace with const
	}

	req := models.Production{
		Name:    name,
		Title:   title,
		Summary: summary,
	}

	resp := models.Production{}
	_, err := cl.post(cl.ns+productionRoute, &req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Productions retrieves a list of productions
func (cl *Client) Productions() (*models.ProductionList, error) {
	if !cl.IsValid() {
		return nil, PodopsClientConfigurationErr
	}

	var resp models.ProductionList
	_, err := cl.get(cl.ns+listProductionsRoute, &resp)
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

	resp := a.StatusObject{}
	status, err := cl.post(cl.ns+fmt.Sprintf(updateResourceRoute, production, kind, guid, force), rsrc, &resp)

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

	status, err := cl.get(cl.ns+fmt.Sprintf(getResourceRoute, production, kind, guid), rsrc)
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

	status, err := cl.get(cl.ns+fmt.Sprintf(findResourceRoute, guid), rsrc)
	if status == http.StatusBadRequest {
		return fmt.Errorf("not found: '%s'", guid)
	}
	if err != nil {
		return err
	}

	return nil
}

// Resources retrieves a list of resources
func (cl *Client) Resources(production, kind string) (*models.ResourceList, error) {
	if !cl.IsValid() {
		return nil, PodopsClientConfigurationErr
	}
	if kind == "" {
		kind = "ALL"
	}

	var resp models.ResourceList
	_, err := cl.get(cl.ns+fmt.Sprintf(listResourcesRoute, production, kind), &resp)
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

	resp := a.StatusObject{}
	status, err := cl.put(cl.ns+fmt.Sprintf(updateResourceRoute, production, kind, guid, force), rsrc, &resp)

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

	status, err := cl.delete(cl.ns+fmt.Sprintf(deleteResourceRoute, production, kind, guid), nil)
	if err != nil {
		return status, err
	}
	return status, nil
}

// Build invokes the BuildEndpoint
func (cl *Client) Build(production string) (*models.BuildRequest, error) {
	if !cl.IsValid() {
		return nil, PodopsClientConfigurationErr
	}

	req := models.BuildRequest{
		GUID: production,
	}
	resp := models.BuildRequest{}

	_, err := cl.post(cl.ns+buildRoute, &req, &resp)
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

	req, err := cl.fileUploadRequest(cl.opts.APIEndpoint+cl.ns+uploadRoute, production, path) // FIXME upload should go against the CDN endpoint
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
