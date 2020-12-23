package podcast

import (
	"net/http"

	"github.com/podops/podops/internal/errors"
	t "github.com/podops/podops/pkg/types"
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
