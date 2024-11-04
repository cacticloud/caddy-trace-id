package caddy_req_id

import (
	"context"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/google/uuid"
)

func init() {
	caddy.RegisterModule(ReqID{})
	httpcaddyfile.RegisterHandlerDirective("req_id", parseCaddyfile)
}

type ReqID struct {
	Enabled bool `json:"enabled,omitempty"`
}

func (ReqID) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.req_id",
		New: func() caddy.Module { return new(ReqID) },
	}
}

func (m *ReqID) Provision(ctx caddy.Context) error {
	return nil
}

func (m ReqID) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	reqID := uuid.New().String()[:32]
	r.Header.Set("Req-ID", reqID)
	w.Header().Set("Req-ID", reqID)

	newContext := context.WithValue(r.Context(), "Req-ID", reqID)
	newRequest := r.WithContext(newContext)

	return next.ServeHTTP(w, newRequest)
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var u ReqID
	if !h.Next() {
		return nil, h.ArgErr()
	}
	remainingArgs := h.RemainingArgs()
	if len(remainingArgs) > 0 {
		u.Enabled = (remainingArgs[0] == "true")
	}
	return u, nil
}

var (
	_ caddy.Provisioner           = (*ReqID)(nil)
	_ caddyhttp.MiddlewareHandler = (*ReqID)(nil)
)
