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
	Logger       *zap.Logger
	CurrentReqID string // 存储为请求处理的一部分的当前 Req-ID
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

func (u *ReqID) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	u.CurrentReqID = uuid.New().String()[:8] // 生成并存储 Req-ID
	r.Header.Set("Req-ID", u.CurrentReqID)

	// 日志中明确记录 Req-ID
	u.Logger.Info("Generated unique ID", zap.String("Req-ID", u.CurrentReqID), zap.String("url", r.URL.String()))

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
