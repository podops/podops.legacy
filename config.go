package podops

import (
	"os/user"
	"path/filepath"

	"github.com/fupas/commons/pkg/env"

	"github.com/podops/podops/internal/cli/netrc"
)

const (
	defaultBucketProduction = "production.podops.dev"
	defaultBucketCDN        = "cdn.podops.dev"

	defaultEndpoint        = "https://podops.dev"
	defaultAPIEndpoint     = "https://api.podops.dev"
	defaultCDNEndpoint     = "https://cdn.podops.dev"
	defaultStorageEndpoint = "https://storage.googleapis.com/cdn.podops.dev"

	machineEntry = "api.podops.dev"
)

var (
	// DefaultEndpoint points to the portal
	DefaultEndpoint string = env.GetString("BASE_URL", defaultEndpoint)

	// DefaultAPIEndpoint points to the API
	DefaultAPIEndpoint string = env.GetString("API_ENDPOINT", defaultAPIEndpoint)

	// DefaultCDNEndpoint points to the CDN
	DefaultCDNEndpoint string = env.GetString("CDN_URL", defaultCDNEndpoint)

	// BucketProduction is the canonical name of the production bucket
	BucketProduction string = env.GetString("BUCKET_PRODUCTION", defaultBucketProduction)

	// BucketCDN is the canonical name of the CDN bucket
	// FIXME will be obsolete
	BucketCDN string = env.GetString("BUCKET_CDN", defaultBucketCDN)

	// StorageEndpoint is the direct link to assets in Google Storage
	StorageEndpoint string = env.GetString("STORAGE_ENDPOINT", defaultStorageEndpoint)
)

// DefaultClientOptions returns a default configuration bases on ENV variables
func DefaultClientOptions() *ClientOption {
	o := ClientOption{
		Token:           env.GetString("PODOPS_API_TOKEN", ""),
		APIEndpoint:     DefaultAPIEndpoint,
		CDNEndpoint:     DefaultCDNEndpoint,
		DefaultEndpoint: DefaultEndpoint,
	}
	return &o
}

func LoadConfiguration() *ClientOption {
	opts := DefaultClientOptions()

	nrc := loadConfig()
	m := nrc.FindMachine(machineEntry)
	if m != nil {
		opts.Token = m.Password
		if m.Account != "" {
			opts.Production = m.Account
		}
	}
	return opts
}

func DefaultConfigPath() string {
	path := env.GetString("PODOPS_CREDENTIALS", "")
	if path == "" {
		usr, _ := user.Current()
		path = filepath.Join(usr.HomeDir, ".netrc")
	}
	return path
}

func loadConfig() *netrc.Netrc {
	nrc, _ := netrc.ParseFile(DefaultConfigPath())
	if nrc == nil {
		nrc = &netrc.Netrc{}
	}
	return nrc
}
