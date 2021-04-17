package cdn

import (
	"net/http"
	"os"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

type (
	ContentHandler struct {
	}
)

var (
	// Interface guards
	_ caddy.Validator             = (*ContentHandler)(nil)
	_ caddy.Provisioner           = (*ContentHandler)(nil)
	_ caddyhttp.MiddlewareHandler = (*ContentHandler)(nil)
	_ caddyfile.Unmarshaler       = (*ContentHandler)(nil)
)

func init() {
	caddy.RegisterModule(ContentHandler{})
	httpcaddyfile.RegisterHandlerDirective("cdn_handler", parseConfig)
}

func (m ContentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	os.Stdout.Write([]byte(r.RequestURI + "\n"))
	return next.ServeHTTP(w, r)
}

func (ContentHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.podops",
		New: func() caddy.Module { return new(ContentHandler) },
	}
}

func (m *ContentHandler) Provision(ctx caddy.Context) error {
	return nil
}

func (m *ContentHandler) Validate() error {
	return nil
}

func (m *ContentHandler) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
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
	var m ContentHandler
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}
