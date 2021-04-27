package podops

import (
	"context"

	"github.com/podops/podops/internal/errordef"
)

// Client is a client for interacting with the PodOps service.
//
// Clients should be reused instead of created as needed.
// The methods of Client are safe for concurrent use by multiple goroutines.
type (
	ClientOption struct {
		Token           string
		Production      string
		APIEndpoint     string
		CDNEndpoint     string
		DefaultEndpoint string
	}

	Client struct {
		opts              *ClientOption
		defaultProduction string
		// internal for now
		realm string
	}
)

// NewClient creates a new podcast client.
//
// Clients should be reused instead of created as needed.
// The methods of a client instance are threadsafe.
func NewClient(ctx context.Context, token string, opts ...*ClientOption) (*Client, error) {

	co := DefaultClientOptions()
	if len(opts) != 0 {
		co = co.Merge(opts[0]) // QA we assume only 1 opts is provided!
	}

	if token != "" {
		co.Token = token
	}

	return New(ctx, co)
}

func New(ctx context.Context, o *ClientOption) (*Client, error) {
	if o == nil || !o.IsValid() {
		return nil, errordef.ErrInvalidClientConfiguration
	}
	return &Client{
		opts:  o,
		realm: "podops",
	}, nil
}

// IsValid checks if all configuration parameters are provided
func (cl *Client) IsValid() bool {
	return cl.opts.IsValid()
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

func (cl *Client) DefaultEndpoint() string {
	return cl.opts.DefaultEndpoint
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
	o := ClientOption{
		Token:           co.Token,
		APIEndpoint:     co.APIEndpoint,
		CDNEndpoint:     co.CDNEndpoint,
		DefaultEndpoint: co.DefaultEndpoint,
	}

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
		if opts.DefaultEndpoint != "" {
			o.DefaultEndpoint = opts.DefaultEndpoint
		}
	}

	return &o
}

// IsValid checks if all configuration parameters are provided
func (co ClientOption) IsValid() bool {
	return co.APIEndpoint != "" && co.CDNEndpoint != "" && co.DefaultEndpoint != ""
}
