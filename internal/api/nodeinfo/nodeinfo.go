package nodeinfo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/api/middleware"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func Router(r *gin.RouterGroup) {
	group := r.Group("/nodeinfo")
	group.Use(middleware.CORS())
	group.GET("/2.0", InfoHandler)
}

func InfoHandler(c *gin.Context) {
	server := c.GetString("proxy-server")
	var err error
	info := models.NodeInfo{
		Version: "2.0",
		Software: models.NodeInfoSoftware{
			Name:    "misstodon",
			Version: global.AppVersion,
		},
		Protocols: []string{"activitypub"},
		Services: models.NodeInfoServices{
			Inbound:  []string{},
			Outbound: []string{},
		},
		Metadata: struct{}{},
	}
	if server != "" {
		info, err = misskey.NodeInfo(
			server,
			models.NodeInfo{
				Version: "2.0",
				Software: models.NodeInfoSoftware{
					Name:    "misstodon",
					Version: global.AppVersion,
				},
				Protocols: []string{"activitypub"},
				Services: models.NodeInfoServices{
					Inbound:  []string{},
					Outbound: []string{},
				},
			})
		if err != nil {
			httperror.AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
	}
	c.JSON(http.StatusOK, info)
}
