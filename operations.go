package podops

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
	productionRoute = "/new"
	// resourceRoute route to call ResourceEndpoint
	resourceRoute = "/update/%s/%s/%s?f=%v" // "/update/:parent/:rsrc/:id"
	// listRoute route to call ListEndpoint
	listRoute = "/list"
	// buildRoute route to call BuildEndpoint
	buildRoute = "/build"
	// uploadRoute route to UploadEndpoint
	uploadRoute = "/upload"
)

// CreateToken creates an access token on the service
// FIXME this is not tested
func (cl *Client) CreateToken(secret, realm, clientID, userID, scope string, duration int64) (string, error) {
	req := a.AuthorizationRequest{
		Secret:     secret,
		Realm:      realm,
		ClientID:   clientID,
		ClientType: "user",
		UserID:     userID,
		Scope:      scope,
		Duration:   duration,
	}
	resp := a.AuthorizationResponse{}

	// create temporary client because we have to swap an existing token with secret
	tempClient, _ := NewClient("")
	tempClient.Token = secret
	status, err := tempClient.post(authenticationRoute, &req, &resp)

	if err != nil {
		return "", fmt.Errorf("create token exception: %v", err)
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("create token exception: %d", status)
	}

	return resp.Token, nil
}

// CreateProduction invokes the CreateProductionEndpoint
func (cl *Client) CreateProduction(name, title, summary string) (*a.ProductionResponse, error) {
	if err := cl.HasToken(); err != nil {
		return nil, err
	}

	if name == "" {
		return nil, fmt.Errorf("resource: name must not be empty")
	}

	req := a.ProductionRequest{
		Name:    name,
		Title:   title,
		Summary: summary,
	}

	resp := a.ProductionResponse{}
	_, err := cl.post(cl.apiNamespace+productionRoute, &req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// List retrieves a list of resources
func (cl *Client) List() (*a.ProductionsResponse, error) {
	if err := cl.HasToken(); err != nil {
		return nil, err
	}

	var resp a.ProductionsResponse
	_, err := cl.get(cl.apiNamespace+listRoute, &resp)
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
	status, err := cl.post(cl.apiNamespace+fmt.Sprintf(resourceRoute, cl.GUID, kind, rsrcGUID, force), rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// UpdateResource invokes the ResourceEndpoint
func (cl *Client) UpdateResource(kind, rsrcGUID string, force bool, rsrc interface{}) (int, error) {
	if err := cl.HasTokenAndGUID(); err != nil {
		return http.StatusBadRequest, err
	}

	resp := a.StatusObject{}
	status, err := cl.put(cl.apiNamespace+fmt.Sprintf(resourceRoute, cl.GUID, kind, rsrcGUID, force), rsrc, &resp)

	if err != nil {
		return status, err
	}
	return status, nil
}

// Build invokes the BuildEndpoint
func (cl *Client) Build(guid string) (string, error) {
	if err := cl.HasTokenAndGUID(); err != nil {
		return "", err
	}

	req := a.BuildRequest{
		GUID: guid,
	}
	resp := a.BuildResponse{}

	_, err := cl.post(cl.apiNamespace+buildRoute, &req, &resp)
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}

// UploadResource invokes the UploadEndpoint
func (cl *Client) UploadResource(path string, force bool) error {
	if err := cl.HasTokenAndGUID(); err != nil {
		return err
	}

	req, err := cl.fileUploadRequest(cl.ServiceEndpoint+cl.apiNamespace+uploadRoute, cl.GUID, path)
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
		return fmt.Errorf("Error uploading '%s': %s", path, resp.Status)
	}

	return nil
}
