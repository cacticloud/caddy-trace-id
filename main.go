package caddy_req_id

import (
	"net/http"
	"strconv"

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
	Logger      *zap.Logger
	LogRequests bool `json:"log_requests"` // 控制是否记录日志
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
	reqID := uuid.New().String()[:8]
	r.Header.Set("Req-ID", reqID)

	if u.LogRequests {
		u.Logger.Info("Generated unique ID", zap.String("Req-ID", reqID), zap.String("url", r.URL.String()))
	}

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
	for h.NextBlock(0) {
		switch h.Val() {
		case "log_requests":
			if !h.NextArg() {
				return nil, h.ArgErr()
			}
			logRequests, err := strconv.ParseBool(h.Val())
			if err != nil {
				return nil, err
			}
			u.LogRequests = logRequests
		default:
			return nil, h.ArgErr()
		}
	}
	return u, nil
}
