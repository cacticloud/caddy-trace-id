package caddy_req_id

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(UniqueID{})
	httpcaddyfile.RegisterHandlerDirective("req_id", parseCaddyfile)
}

type UniqueID struct {
	UserIDHeader string `json:"user_id_header,omitempty"` // 用户ID从哪个HTTP头部字段读取
	Logger       *zap.Logger
}

func (UniqueID) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.req_id",
		New: func() caddy.Module { return new(UniqueID) },
	}
}

// Provision set up the module's configuration
func (u *UniqueID) Provision(ctx caddy.Context) error {
	u.Logger = ctx.Logger(u) // Get the logger from the context
	return nil
}

func (u UniqueID) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	userID := r.Header.Get(u.UserIDHeader)
	if userID == "" {
		userID = "anonymous"
	}

	uniqueID := generateUniqueID(time.Now(), r.RemoteAddr, userID, r.Method, r.URL.String())
	r.Header.Set("X-Unique-ID", uniqueID) // 设置请求头，传递到后端

	newContext := context.WithValue(r.Context(), "X-Unique-ID", uniqueID)
	newRequest := r.WithContext(newContext)

	return next.ServeHTTP(w, newRequest)
}

func generateUniqueID(t time.Time, ip, userID, method, url string) string {
	data := fmt.Sprintf("%v-%v-%v-%v-%v", t.UnixNano(), ip, userID, method, url)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

var (
	_ caddyhttp.MiddlewareHandler = (*UniqueID)(nil)
)

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
