package modules

import (
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/txsvc/platform/v2"
)

type (
	StorageModuleImpl struct {
	}
)

var (
	// Interface guards
	_ caddy.Validator             = (*StorageModuleImpl)(nil)
	_ caddy.Provisioner           = (*StorageModuleImpl)(nil)
	_ caddyhttp.MiddlewareHandler = (*StorageModuleImpl)(nil)
	_ caddyfile.Unmarshaler       = (*StorageModuleImpl)(nil)
)

func init() {
	caddy.RegisterModule(StorageModuleImpl{})
	httpcaddyfile.RegisterHandlerDirective("cdn_storage", parseConfig)
}

func (m StorageModuleImpl) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	parts := strings.Split(r.RequestURI, "/")
	if len(parts) > 1 {
		prod := parts[0]
		asset := parts[1]
		contentType := "unknown"
		size := "0"

		// track api access for billing etc
		platform.Meter(platform.NewHttpContext(r), "cdn.storage", "production", prod, "user-agent", r.UserAgent(), "remote_addr", r.RemoteAddr, "type", contentType, "name", asset, "size", size)
	}

	//os.Stdout.Write([]byte(r.RequestURI + "\n"))
	return next.ServeHTTP(w, r)
}

func (StorageModuleImpl) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.podops",
		New: func() caddy.Module { return new(StorageModuleImpl) },
	}
}

func (m *StorageModuleImpl) Provision(ctx caddy.Context) error {
	return nil
}

func (m *StorageModuleImpl) Validate() error {
	return nil
}

func (m *StorageModuleImpl) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	/*
		for d.Next() {
			if !d.Args(&m.Output) {
				return d.ArgErr()
			}
		}
	*/
	return nil
}

func parseConfig(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m StorageModuleImpl
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}
