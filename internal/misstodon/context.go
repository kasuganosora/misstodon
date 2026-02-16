package misstodon

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/utils"
)

type Context struct {
	m sync.Map
}

func ContextWithGinContext(gCtx *gin.Context, tokenRequired ...bool) (*Context, error) {
	c := &Context{}
	if server, ok := gCtx.Get("proxy-server"); ok {
		c.SetProxyServer(server.(string))
	}
	if len(tokenRequired) > 0 && tokenRequired[0] {
		token, err := utils.GetHeaderToken(gCtx.Request.Header)
		if err != nil && (len(tokenRequired) > 0 && tokenRequired[0]) {
			return nil, gin.Error{Err: err, Type: gin.ErrorTypePublic, Meta: gin.H{"code": http.StatusUnauthorized, "message": "the access token is invalid"}}
		}
		arr := strings.Split(token, ".")
		if len(arr) < 2 {
			return nil, gin.Error{Err: err, Type: gin.ErrorTypePublic, Meta: gin.H{"code": http.StatusUnauthorized, "message": "the access token is invalid"}}
		}
		c.SetUserID(arr[0])
		c.SetToken(arr[1])
	}
	c.SetHOST(gCtx.Request.Host)
	return c, nil
}

func ContextWithValues(proxyServer, token string) *Context {
	c := &Context{}
	c.SetProxyServer(proxyServer)
	c.SetToken(token)
	return c
}

func (*Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*Context) Done() <-chan struct{} {
	return nil
}

func (*Context) Err() error {
	return nil
}

func (c *Context) Value(key any) any {
	if val, ok := c.m.Load(key); ok {
		return val
	}
	return nil
}

func (c *Context) SetValue(key any, val any) {
	c.m.Store(key, val)
}

func (c *Context) String(key string) *string {
	if val, ok := c.m.Load(key); ok {
		valStr := val.(string)
		return &valStr
	}
	return nil
}

func (c *Context) ProxyServer() string {
	if s := c.String("proxy-server"); s != nil {
		return *s
	}
	return ""
}

func (c *Context) SetProxyServer(val string) {
	c.SetValue("proxy-server", val)
}

func (c *Context) Token() *string {
	return c.String("token")
}

func (c *Context) SetToken(val string) {
	c.SetValue("token", val)
}

func (c *Context) UserID() *string {
	return c.String("user_id")
}

func (c *Context) SetUserID(val string) {
	c.SetValue("user_id", val)
}

func (c *Context) HOST() *string {
	return c.String("host")
}

func (c *Context) SetHOST(val string) {
	c.SetValue("host", val)
}
