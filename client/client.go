package client

import (
	"context"
	"fmt"

	"github.com/podops/podops/apiv1"
)

var (
	PodopsClientConfigurationErr error = fmt.Errorf("new client: invalid configuration")
)

// Client is a client for interacting with the PodOps service.
//
// Clients should be reused instead of created as needed.
// The methods of Client are safe for concurrent use by multiple goroutines.
type (
	ClientOption struct {
		Token          string
		Production     string
		APIEndpoint    string
		CDNEndpoint    string
		PortalEndpoint string
	}

	Client struct {
		opts              *ClientOption
		defaultProduction string
		// internal for now
		validated bool
		valid     bool
		ns        string
		realm     string
	}
)

// IsValid checks if all configuration parameters are provided
func (cl *Client) IsValid() bool {
	if cl.validated {
		return cl.valid
	}
	cl.validated = true
	// verify the opts first
	if !cl.opts.IsValid() {
		cl.valid = false
		return false
	}

	cl.valid = true // FIXME try to verify the token against the API

	return true
}

func (cl *Client) SetProduction(production string) {
	cl.defaultProduction = production
}

func (cl *Client) APIEndpoint() string {
	return cl.opts.APIEndpoint
}

func (cl *Client) CDNEndpoint() string {
	return cl.opts.CDNEndpoint
}

func (cl *Client) PortalEndpoint() string {
	return cl.opts.PortalEndpoint
}

func (cl *Client) Realm() string {
	return cl.realm
}

func (cl *Client) Token() string {
	return cl.opts.Token
}

func (cl *Client) DefaultProduction() string {
	return cl.defaultProduction
}

// Merge clones co and combines it with the provided options
func (co ClientOption) Merge(opts *ClientOption) *ClientOption {
	o := &ClientOption{}
	o.Token = co.Token
	o.APIEndpoint = co.APIEndpoint
	o.CDNEndpoint = co.CDNEndpoint
	o.PortalEndpoint = co.PortalEndpoint

	if opts != nil {
		if opts.Token != "" {
			o.Token = opts.Token
		}
		if opts.APIEndpoint != "" {
			o.APIEndpoint = opts.APIEndpoint
		}
		if opts.CDNEndpoint != "" {
			o.CDNEndpoint = opts.CDNEndpoint
		}
		if opts.PortalEndpoint != "" {
			o.PortalEndpoint = opts.PortalEndpoint
		}
	}

	return o
}

// IsValid checks if all configuration parameters are provided
func (co ClientOption) IsValid() bool {
	return co.APIEndpoint != "" && co.CDNEndpoint != "" && co.PortalEndpoint != ""
	// we can not validate token and production. There are some API calls that do not need them...
}

func New(ctx context.Context, o *ClientOption) (*Client, error) {
	if o == nil || !o.IsValid() {
		return nil, PodopsClientConfigurationErr
	}
	return &Client{
		opts:      o,
		validated: false,
		valid:     false,
		ns:        apiv1.NamespacePrefix,
		realm:     "podops",
	}, nil
}
