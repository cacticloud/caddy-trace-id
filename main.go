package caddytraceid

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Gizmo{})
}

type Gizmo struct{}

// CaddyModule 返回 Caddy 模块的信息，现在使用 http.handlers 命名空间
func (Gizmo) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.gizmo", // 修改这里以匹配 http.handlers
		New: func() caddy.Module { return new(Gizmo) },
	}
}

// ServeHTTP 实现 caddyhttp.MiddlewareHandler 接口
func (g *Gizmo) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// 实现你的处理逻辑
	return next.ServeHTTP(w, r)
}

var (
	_ caddyhttp.MiddlewareHandler = (*Gizmo)(nil) // 确保 Gizmo 实现了 MiddlewareHandler 接口
)
