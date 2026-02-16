package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey/streaming"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var wsUpgrade = websocket.Upgrader{
	ReadBufferSize:  4096, // we don't expect reads
	WriteBufferSize: 4096,
	Subprotocols:    []string{},
	CheckOrigin: func(r *http.Request) bool { return true },
}

func StreamingRouter(r *gin.RouterGroup) {
	r.GET("/streaming", StreamingHandler)
}

func StreamingHandler(c *gin.Context) {
	var token string
	if token = c.Query("access_token"); token == "" {
		if token = c.Request.Header.Get("Sec-Websocket-Protocol"); token == "" {
			httperror.AbortWithError(c, http.StatusBadRequest, errors.New("no access token provided"))
			return
		}
	}
	server := c.GetString("proxy-server")

	conn, err := wsUpgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan models.StreamEvent)
	defer close(ch)
	go func() {
		if err := streaming.Streaming(ctx, server, token, ch); err != nil {
			log.Debug().Caller().Err(err).Msg("Streaming error")
		}
		_ = conn.Close()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-ch:
				log.Debug().Caller().Any("event", event).Msg("Streaming")
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		_, _, err = conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				cancel()
				return
			}
			httperror.AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
	}
}
