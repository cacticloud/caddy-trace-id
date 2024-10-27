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
	caddy.RegisterModule(UniqueID{})
	httpcaddyfile.RegisterHandlerDirective("unique_id", parseCaddyfile)
}

type UniqueID struct {
	UserIDHeader string `json:"user_id_header,omitempty"`
	Logger       *zap.Logger
	Prefix       string `json:"prefix,omitempty"`
}

func (UniqueID) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.unique_id",
		New: func() caddy.Module { return new(UniqueID) },
	}
}

func (u *UniqueID) Provision(ctx caddy.Context) error {
	u.Logger = ctx.Logger(u)
	if u.Prefix == "" {
		u.Prefix = "req" // 设置默认前缀，如果未在配置中指定
	}
	return nil
}

func (u UniqueID) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	partUUID := uuid.New().String()[:8]
	uniqueID := u.Prefix + "-" + partUUID
	r.Header.Set("Req-ID", uniqueID)

	u.Logger.Info("Generated unique ID", zap.String("Req-ID", uniqueID), zap.String("url", r.URL.String()))

	return next.ServeHTTP(w, r)
}

var (
	_ caddyhttp.MiddlewareHandler = (*UniqueID)(nil)
)

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var u UniqueID
	if !h.Next() {
		return nil, h.ArgErr()
	}
	if h.Args(&u.UserIDHeader, &u.Prefix) {
		return u, nil
	}
	return nil, h.ArgErr()
}
