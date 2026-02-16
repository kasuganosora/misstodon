package httperror

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

type ServerError struct {
	TraceID string `json:"trace_id,omitempty"`
	Error   string `json:"error"`
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last()
		code := http.StatusInternalServerError
		info := ServerError{Error: err.Error()}

		if code == http.StatusInternalServerError {
			id := xid.New().String()
			info = ServerError{
				TraceID: id,
				Error:   "Internal Server Error",
			}
			log.Warn().Err(err).
				Str("user_agent", c.Request.UserAgent()).
				Str("trace_id", id).
				Int("code", code).
				Msg("Server Error")
		}
		c.JSON(code, info)
	}
}

func AbortWithError(c *gin.Context, code int, err error) {
	c.Abort()
	info := ServerError{Error: err.Error()}
	if code == http.StatusInternalServerError {
		id := xid.New().String()
		info = ServerError{
			TraceID: id,
			Error:   "Internal Server Error",
		}
		log.Warn().Err(err).
			Str("user_agent", c.Request.UserAgent()).
			Str("trace_id", id).
			Int("code", code).
			Msg("Server Error")
	}
	c.JSON(code, info)
}
