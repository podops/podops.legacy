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

// SetProduction sets the context of further operations
func (cl *Client) SetProduction(guid string) {
	cl.GUID = guid
	// FIXME make sure we own the GUID
}

// CreateToken creates an access token on the server
// FIXME this is not tested
func (cl *Client) CreateToken(secret, realm, clientID, userID, scope string, duration int64) (string, error) {
	req := a.OldAuthorizationRequest{
		Secret:     secret,
		Realm:      realm,
		ClientID:   clientID,
		ClientType: "user",
		UserID:     userID,
		Scope:      scope,
		Duration:   duration,
	}
	resp := a.OldAuthorizationResponse{}

	// create temporary client because we have to swap an existing token with secret
	tempClient := DefaultClient("")
	tempClient.Token = secret
	status, err := tempClient.post(authenticationRoute, &req, &resp)

	if err != nil {
		return "", fmt.Errorf("create token exception: %v", err)
	}
	if status != http.StatusCreated {
		return "", fmt.Errorf("create token exception: %d", status)
	}

	return resp.Token, nil
}

// CreateProduction invokes the CreateProductionEndpoint
func (cl *Client) CreateProduction(name, title, summary string) (*a.Production, error) {
	if err := cl.HasToken(); err != nil {
		return nil, err
	}

	if name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}

	req := a.Production{
		Name:    name,
		Title:   title,
		Summary: summary,
	}

	resp := a.Production{}
	_, err := cl.post(cl.Namespace+productionRoute, &req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Productions retrieves a list of productions
func (cl *Client) Productions() (*a.ProductionList, error) {
	if err := cl.HasToken(); err != nil {
		return nil, err
	}

	var resp a.ProductionList
	_, err := cl.get(cl.Namespace+listProductionsRoute, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateResource invokes the ResourceEndpoint
func (cl *Client) CreateResource(kind, rsrcGUID string, force bool, rsrc interface{}) (int, error) {
	if err := cl.HasTokenAndGUID(); err != nil {
		return http.StatusBadRequest, err
	}

	resp := a.StatusObject{}
	status, err := cl.post(cl.Namespace+fmt.Sprintf(updateResourceRoute, cl.GUID, kind, rsrcGUID, force), rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// GetResource returns a resource file
func (cl *Client) GetResource(prod, kind, guid string, rsrc interface{}) error {
	if err := cl.HasToken(); err != nil {
		return err
	}

	status, err := cl.get(cl.Namespace+fmt.Sprintf(getResourceRoute, prod, kind, guid), rsrc)
	if status == http.StatusBadRequest {
		return fmt.Errorf("not found: '%s/%s-%s'", prod, kind, guid)
	}
	if err != nil {
		return err
	}

	return nil
}

// UpdateResource invokes the ResourceEndpoint
func (cl *Client) UpdateResource(kind, rsrcGUID string, force bool, rsrc interface{}) (int, error) {
	if err := cl.HasTokenAndGUID(); err != nil {
		return http.StatusBadRequest, err
	}

	resp := a.StatusObject{}
	status, err := cl.put(cl.Namespace+fmt.Sprintf(updateResourceRoute, cl.GUID, kind, rsrcGUID, force), rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// Resources retrieves a list of resources
func (cl *Client) Resources(prod, kind string) (*a.ResourceList, error) {
	if err := cl.HasToken(); err != nil {
		return nil, err
	}
	if kind == "" {
		kind = "ALL"
	}

	var resp a.ResourceList
	_, err := cl.get(cl.Namespace+fmt.Sprintf(listResourcesRoute, prod, kind), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteResource deletes a resources
func (cl *Client) DeleteResource(prod, kind, guid string) (int, error) {
	if err := cl.HasToken(); err != nil {
		return http.StatusBadRequest, err
	}

	status, err := cl.delete(cl.Namespace+fmt.Sprintf(deleteResourceRoute, prod, kind, guid), nil)
	if err != nil {
		return status, err
	}
	return status, nil
}

// Build invokes the BuildEndpoint
func (cl *Client) Build(guid string) (*a.Build, error) {
	if err := cl.HasTokenAndGUID(); err != nil {
		return nil, err
	}

	req := a.Build{
		GUID: guid,
	}
	resp := a.Build{}

	_, err := cl.post(cl.Namespace+buildRoute, &req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Upload invokes the UploadEndpoint
func (cl *Client) Upload(path string, force bool) error {
	if err := cl.HasTokenAndGUID(); err != nil {
		return err
	}

	req, err := cl.fileUploadRequest(cl.ServiceEndpoint+cl.Namespace+uploadRoute, cl.GUID, path)
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
