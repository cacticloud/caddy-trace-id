package caddy_unique_id

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(UniqueID{})
	httpcaddyfile.RegisterHandlerDirective("unique_id", parseCaddyfile)
}

type UniqueID struct {
	UserIDHeader string `json:"user_id_header,omitempty"` // 用户ID从哪个HTTP头部字段读取
}

func (UniqueID) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.unique_id",
		New: func() caddy.Module { return new(UniqueID) },
	}
}

func (u UniqueID) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	userID := r.Header.Get(u.UserIDHeader)
	if userID == "" {
		userID = "anonymous"
	}

	// 生成唯一 ID
	uniqueID := generateUniqueID(time.Now(), r.RemoteAddr, userID)
	r.Header.Set("X-Unique-ID", uniqueID) // 将此 ID 加入请求头中，以便后续使用

	// 可选：将此 ID 加入响应头或日志
	w.Header().Set("X-Unique-ID", uniqueID)

	return next.ServeHTTP(w, r)
}

// generateUniqueID 根据时间、IP 和用户ID生成一个哈希值
func generateUniqueID(t time.Time, ip, userID string) string {
	data := fmt.Sprintf("%v-%v-%v", t.UnixNano(), ip, userID)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

var (
	_ caddyhttp.MiddlewareHandler = (*UniqueID)(nil)
)

// parseCaddyfile 用于解析 Caddyfile 并配置 UniqueID 实例
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var u UniqueID
	if !h.Next() {
		return nil, h.ArgErr()
	}
	if !h.AllArgs(&u.UserIDHeader) {
		return nil, h.ArgErr()
	}
	return u, nil
}
