package podops

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/txsvc/platform/pkg/api"

	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/messagedef"
	"github.com/podops/podops/internal/transport"
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
	// uploadRoute route to the CDN UploadEndpoint
	uploadRoute = "/_w/upload"
)

func assertNotEmpty(claims ...string) bool {
	if len(claims) == 0 {
		return false
	}
	for _, s := range claims {
		if s == "" {
			return false
		}
	}
	return true
}

// CreateProduction invokes the CreateProductionEndpoint
func (cl *Client) CreateProduction(name, title, summary string) (*Production, error) {
	if !cl.IsValid() {
		return nil, errordef.ErrInvalidClientConfiguration
	}
	if name == "" {
		return nil, errordef.ErrInvalidParameters
	}

	req := Production{
		Name:    name,
		Title:   title,
		Summary: summary,
	}

	resp := Production{}
	_, err := transport.Post(cl.opts.APIEndpoint, productionRoute, cl.opts.Token, &req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Productions retrieves a list of productions
func (cl *Client) Productions() (*ProductionList, error) {
	if !cl.IsValid() {
		return nil, errordef.ErrInvalidClientConfiguration
	}

	var resp ProductionList
	_, err := transport.Get(cl.opts.APIEndpoint, listProductionsRoute, cl.opts.Token, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateResource invokes the ResourceEndpoint
func (cl *Client) CreateResource(production, kind, guid string, force bool, rsrc interface{}) (int, error) {
	if !cl.IsValid() {
		return http.StatusBadRequest, errordef.ErrInvalidClientConfiguration
	}
	if !assertNotEmpty(production, kind, guid) {
		return http.StatusBadRequest, errordef.ErrInvalidParameters
	}

	resp := api.StatusObject{}
	status, err := transport.Post(cl.opts.APIEndpoint, fmt.Sprintf(updateResourceRoute, production, kind, guid, force), cl.opts.Token, rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// GetResource returns a resource file
func (cl *Client) GetResource(production, kind, guid string, rsrc interface{}) error {
	if !cl.IsValid() {
		return errordef.ErrInvalidClientConfiguration
	}
	if !assertNotEmpty(production, kind, guid) {
		return errordef.ErrInvalidParameters
	}

	status, err := transport.Get(cl.opts.APIEndpoint, fmt.Sprintf(getResourceRoute, production, kind, guid), cl.opts.Token, rsrc)
	if status == http.StatusBadRequest {
		return fmt.Errorf(messagedef.MsgResourceNotFound, fmt.Sprintf("%s/%s-%s", production, kind, guid))
	}
	if err != nil {
		return err
	}

	return nil
}

// FindResource returns a resource file
func (cl *Client) FindResource(guid string, rsrc interface{}) error {
	if !cl.IsValid() {
		return errordef.ErrInvalidClientConfiguration
	}
	if guid == "" {
		return errordef.ErrInvalidParameters
	}

	status, err := transport.Get(cl.opts.APIEndpoint, fmt.Sprintf(findResourceRoute, guid), cl.opts.Token, rsrc)
	if status == http.StatusBadRequest {
		return fmt.Errorf(messagedef.MsgResourceNotFound, guid)
	}
	if err != nil {
		return err
	}

	return nil
}

// Resources retrieves a list of resources
func (cl *Client) Resources(production, kind string) (*ResourceList, error) {
	if !cl.IsValid() {
		return nil, errordef.ErrInvalidClientConfiguration
	}
	if production == "" {
		return nil, errordef.ErrInvalidParameters
	}
	if kind == "" {
		kind = "ALL"
	}

	var resp ResourceList
	_, err := transport.Get(cl.opts.APIEndpoint, fmt.Sprintf(listResourcesRoute, production, kind), cl.opts.Token, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateResource invokes the ResourceEndpoint
func (cl *Client) UpdateResource(production, kind, guid string, force bool, rsrc interface{}) (int, error) {
	if !cl.IsValid() {
		return http.StatusBadRequest, errordef.ErrInvalidClientConfiguration
	}
	if !assertNotEmpty(production, kind, guid) {
		return http.StatusBadRequest, errordef.ErrInvalidParameters
	}

	resp := api.StatusObject{}
	status, err := transport.Put(cl.opts.APIEndpoint, fmt.Sprintf(updateResourceRoute, production, kind, guid, force), cl.opts.Token, rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// DeleteResource deletes a resources
func (cl *Client) DeleteResource(production, kind, guid string) (int, error) {
	if !cl.IsValid() {
		return http.StatusBadRequest, errordef.ErrInvalidClientConfiguration
	}
	if !assertNotEmpty(production, kind, guid) {
		return http.StatusBadRequest, errordef.ErrInvalidParameters
	}

	status, err := transport.Delete(cl.opts.APIEndpoint, fmt.Sprintf(deleteResourceRoute, production, kind, guid), cl.opts.Token, nil)
	if err != nil {
		return status, err
	}
	return status, nil
}

// Build invokes the BuildEndpoint
func (cl *Client) Build(production string) (*BuildRequest, error) {
	if !cl.IsValid() {
		return nil, errordef.ErrInvalidClientConfiguration
	}
	if production == "" {
		return nil, errordef.ErrInvalidParameters
	}

	req := BuildRequest{
		GUID: production,
	}
	resp := BuildRequest{}

	_, err := transport.Post(cl.opts.APIEndpoint, buildRoute, cl.opts.Token, &req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Upload invokes the UploadEndpoint
func (cl *Client) Upload(production, path string, force bool) error {
	if !cl.IsValid() {
		return errordef.ErrInvalidClientConfiguration
	}
	if !assertNotEmpty(production, path) {
		return errordef.ErrInvalidParameters
	}

	req, err := transport.Upload(cl.opts.CDNEndpoint, uploadRoute, cl.opts.Token, production, "asset", path)
	if err != nil {
		log.Fatal(err)
	}

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
		return fmt.Errorf(messagedef.MsgResourceUploadError, fmt.Sprintf("%s:%s", path, resp.Status))
	}

	return nil
}
