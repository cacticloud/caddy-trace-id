package caddy_extra

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(ExtraInfo{})
	httpcaddyfile.RegisterHandlerDirective("extra_info", parseCaddyfile)
}

// ExtraInfo 是一个简单的 HTTP 中间件，用于在响应中添加额外信息
type ExtraInfo struct {
	Message string `json:"message,omitempty"`
}

// CaddyModule 返回一个关于插件的信息
func (ExtraInfo) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.extra_info",
		New: func() caddy.Module { return new(ExtraInfo) },
	}
}

// ServeHTTP 实现 caddyhttp.MiddlewareHandler 接口
func (e ExtraInfo) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// 在响应头中设置额外信息
	w.Header().Add("X-Extra-Info", e.Message)
	return next.ServeHTTP(w, r) // 继续执行下一个处理器
}

var (
	_ caddyhttp.MiddlewareHandler = (*ExtraInfo)(nil)
)

// parseCaddyfile 用于解析 Caddyfile 并配置 ExtraInfo 实例
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var e ExtraInfo
	if !h.Next() {
		return nil, h.ArgErr()
	}
	if !h.AllArgs(&e.Message) {
		return nil, h.ArgErr()
	}
	return e, nil
}
