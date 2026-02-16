package wellknown

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/api/middleware"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func Router(r *gin.RouterGroup) {
	group := r.Group("/.well-known")
	group.Use(middleware.CORS())
	group.GET("/nodeinfo", NodeInfoHandler)
	group.GET("/webfinger", WebFingerHandler)
	group.GET("/host-meta", HostMeta)
}

func NodeInfoHandler(c *gin.Context) {
	server := c.GetString("proxy-server")
	href := "https://" + c.Request.Host + "/nodeinfo/2.0"
	if server != "" {
		href += "?server=" + server
	}
	c.JSON(http.StatusOK, utils.Map{
		"links": []utils.StrMap{
			{
				"rel":  "http://nodeinfo.diaspora.software/ns/schema/2.0",
				"href": href,
			},
		},
	})
}

func WebFingerHandler(c *gin.Context) {
	resource := c.Query("resource")
	if resource == "" {
		c.JSON(http.StatusBadRequest, httperror.ServerError{
			Error: "resource is required",
		})
		return
	}
	misskey.WebFinger(c.GetString("proxy-server"), resource, c.Writer)
}

func HostMeta(c *gin.Context) {
	misskey.HostMeta(c.GetString("proxy-server"), c.Writer)
}
