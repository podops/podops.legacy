package podops

import (
	"context"

	"github.com/fupas/commons/pkg/env"

	a "github.com/podops/podops/apiv1"
	cl "github.com/podops/podops/pkg/client"
)

// NewClient creates a new podcast client.
//
// Clients should be reused instead of created as needed.
// The methods of a client instance are threadsafe.
func NewClient(ctx context.Context, token string, opts ...*cl.ClientOption) (*cl.Client, error) {

	co := DefaultClientOptions()
	if len(opts) != 0 {
		// FIXME we assume only 1 opts is provided!
		co = co.Merge(opts[0])
	}

	if token != "" {
		co.Token = token
	}

	return cl.New(ctx, co)
}

// DefaultClientOptions returns a default client configuration bases on ENV variables
func DefaultClientOptions() *cl.ClientOption {
	o := &cl.ClientOption{}
	o.Token = env.GetString("PODOPS_API_TOKEN", "")
	o.APIEndpoint = a.DefaultAPIEndpoint
	o.CDNEndpoint = a.DefaultCDNEndpoint
	o.PortalEndpoint = a.DefaultPortalEndpoint
	return o
}
