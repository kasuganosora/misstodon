package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Logger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	c.Next()

	latency := time.Since(start)
	status := c.Writer.Status()
	method := c.Request.Method
	proxyServer, _ := c.Get("proxy-server")

	if raw != "" {
		path = path + "?" + raw
	}

	log.Info().
		Str("proxy-server", proxyServer.(string)).
		Str("method", method).
		Str("uri", path).
		Int("status", status).
		Str("latency", latency.String()).
		Msg("Request")
}
