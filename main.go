package caddy_req_id

import (
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
	Logger *zap.Logger
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
	ReqID := uuid.New().String()[:8]
	r.Header.Set("Req-ID", ReqID)

	u.Logger.Info("Generated unique ID", zap.String("Req-ID", ReqID), zap.String("url", r.URL.String()))

	return next.ServeHTTP(w, r)
}

var (
	_ caddyhttp.MiddlewareHandler = (*ReqID)(nil)
)

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var u ReqID
	if !h.Next() {
		return nil, h.ArgErr()
	}
	return u, nil
}
