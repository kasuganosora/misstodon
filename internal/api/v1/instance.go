package v1

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
)

func InstanceRouter(e *echo.Group) {
	group := e.Group("/instance")
	group.GET("", InstanceHandler)
	group.GET("/peers", InstancePeersHandler)
	e.GET("/custom_emojis", InstanceCustomEmojis)
}

func InstanceHandler(c echo.Context) error {
	info, err := misskey.Instance(
		c.Get("proxy-server").(string),
		global.AppVersion)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, info)
}

func InstancePeersHandler(c echo.Context) error {
	peers, err := misskey.InstancePeers(c.Get("proxy-server").(string))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, utils.SliceIfNull(peers))
}

func InstanceCustomEmojis(c echo.Context) error {
	server := c.Get("proxy-server").(string)
	emojis, err := misskey.InstanceCustomEmojis(server)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, emojis)
}
