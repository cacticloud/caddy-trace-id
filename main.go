package caddy_req_id

import (
	"context"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(ReqID{})
	httpcaddyfile.RegisterHandlerDirective("req_id", parseCaddyfile)
}

type ReqID struct {
	Enabled bool `json:"enabled,omitempty"`
	Logger  *zap.Logger
}

func (ReqID) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.req_id",
		New: func() caddy.Module { return new(ReqID) },
	}
}

func (u *ReqID) Provision(ctx caddy.Context) error {
	u.Logger = ctx.Logger(u)
	return nil
}

func (u ReqID) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	if !u.Enabled {
		return next.ServeHTTP(w, r)
	}

	reqID := uuid.New().String()[:32]
	r.Header.Set("Req-ID", reqID)
	w.Header().Set("Req-ID", reqID)

	newContext := context.WithValue(r.Context(), "Req-ID", reqID)
	newRequest := r.WithContext(newContext)

	return next.ServeHTTP(w, newRequest)
}

var (
	_ caddyhttp.MiddlewareHandler = (*ReqID)(nil)
)

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var u ReqID
	if !h.Next() {
		return nil, h.ArgErr()
	}
	if h.Args(&u.Enabled) {
		return u, nil
	}
	u.Enabled = true
	return u, nil
}
