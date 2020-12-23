package podcast

import (
	t "github.com/podops/podops/pkg/types"
)

// CreateProduction invokes the CreateProductionEndpoint
func (cl *Client) CreateProduction(name, title, summary string) (*t.ProductionResponse, error) {

	// FIXME param validation

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
