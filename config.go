package podops

import (
	"os/user"
	"path/filepath"

	"github.com/txsvc/platform/v2/pkg/env"
	"github.com/txsvc/platform/v2/pkg/netrc"
)

const (
	defaultBucketProduction = "production.podops.dev"

	defaultEndpoint        = "https://podops.dev"
	defaultAPIEndpoint     = "https://api.podops.dev"
	defaultCDNEndpoint     = "https://cdn.podops.dev"
	defaultStorageEndpoint = "https://storage.podops.dev"
	defaultStorageLocation = "/data/storage/cdn"

	machineEntry = "api.podops.dev"
)

var (
	// DefaultEndpoint points to the portal
	DefaultEndpoint string = env.GetString("BASE_URL", defaultEndpoint)

	// DefaultAPIEndpoint points to the API
	DefaultAPIEndpoint string = env.GetString("API_ENDPOINT", defaultAPIEndpoint)

	// DefaultCDNEndpoint points to the CDN
	DefaultCDNEndpoint string = env.GetString("CDN_URL", defaultCDNEndpoint)

	// DefaultStorageEndpoint is the direct link to assets in the CDN
	DefaultStorageEndpoint string = env.GetString("STORAGE_ENDPOINT", defaultStorageEndpoint)

	// BucketProduction is the canonical name of the production bucket
	BucketProduction string = env.GetString("BUCKET_PRODUCTION", defaultBucketProduction)

	// StorageLocation is the root location for the cdn
	StorageLocation = env.GetString("STORAGE_LOCATION", defaultStorageLocation)
)

// DefaultClientOptions returns a default configuration bases on ENV variables
func DefaultClientOptions() *ClientOption {
	o := ClientOption{
		Token:           env.GetString("PODOPS_API_KEY", ""),
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
