package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func InstanceRouter(r *gin.RouterGroup) {
	group := r.Group("/instance")
	group.GET("", InstanceHandler)
	group.GET("/peers", InstancePeersHandler)
	group.GET("/rules", InstanceRulesHandler)
	r.GET("/custom_emojis", InstanceCustomEmojis)
}

func InstanceHandler(c *gin.Context) {
	info, err := misskey.Instance(
		c.GetString("proxy-server"),
		global.AppVersion,
		c.Request.Host)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func InstanceRulesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, []any{})
}

func InstancePeersHandler(c *gin.Context) {
	peers, err := misskey.InstancePeers(c.GetString("proxy-server"))
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(peers))
}

func InstanceCustomEmojis(c *gin.Context) {
	server := c.GetString("proxy-server")
	emojis, err := misskey.InstanceCustomEmojis(server)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, emojis)
}
