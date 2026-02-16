package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/middleware"
	"github.com/gizmo-ds/misstodon/internal/api/nodeinfo"
	"github.com/gizmo-ds/misstodon/internal/api/oauth"
	v1 "github.com/gizmo-ds/misstodon/internal/api/v1"
	v2 "github.com/gizmo-ds/misstodon/internal/api/v2"
	"github.com/gizmo-ds/misstodon/internal/api/wellknown"
	"github.com/gizmo-ds/misstodon/internal/global"
)

func Router(r *gin.Engine) {
	r.Use(middleware.SetContextData)
	r.Use(gin.Recovery())
	if global.Config.Logger.RequestLogger {
		r.Use(middleware.Logger)
	}

	for _, group := range []*gin.RouterGroup{
		r.Group(""),
		r.Group("/:proxyServer"),
	} {
		wellknown.Router(group)
		nodeinfo.Router(group)
		oauth.Router(group)
		v1Api := group.Group("/api/v1")
		v1Api.Use(middleware.CORS())
		v2Api := group.Group("/api/v2")
		v2Api.Use(middleware.CORS())
		group.GET("/static/missing.png", v1.MissingImageHandler)
		v1.InstanceRouter(v1Api)
		v1.AccountsRouter(v1Api)
		v1.ApplicationRouter(v1Api)
		v1.StatusesRouter(v1Api)
		v1.StreamingRouter(v1Api)
		v1.TimelinesRouter(v1Api)
		v1.TrendsRouter(v1Api)
		v1.MediaRouter(v1Api)
		v1.NotificationsRouter(v1Api)
		v1.PollsRouter(v1Api)
		v1.BlocksRouter(v1Api)
		v1.MutesRouter(v1Api)
		v1.ReportsRouter(v1Api)
		v1.AnnouncementsRouter(v1Api)
		v2.MediaRouter(v2Api)
		v2.SearchRouter(v2Api)
		v2.InstanceRouter(v2Api)
		v2.SuggestionsRouter(v2Api)

		v1Api.GET("/bookmarks", v1.StatusBookmarks)
		v1Api.GET("/follow_requests", v1.AccountFollowRequests)
		v1Api.POST("/follow_requests/:id/authorize", v1.FollowRequestAuthorize)
		v1Api.POST("/follow_requests/:id/reject", v1.FollowRequestReject)
		v1Api.GET("/suggestions", v1.SuggestionsHandler)
		v1Api.GET("/preferences", v1.PreferencesHandler)
		v1Api.GET("/markers", v1.MarkersGetHandler)
		v1Api.POST("/markers", v1.MarkersPostHandler)
		v1Api.GET("/conversations", v1.ConversationsHandler)
		v1Api.GET("/followed_tags", func(c *gin.Context) { c.JSON(200, []any{}) })
		v1Api.GET("/endorsements", func(c *gin.Context) { c.JSON(200, []any{}) })
		v1Api.GET("/lists", func(c *gin.Context) { c.JSON(200, []any{}) })
		v1Api.GET("/domain_blocks", func(c *gin.Context) { c.JSON(200, []any{}) })
		v1Api.GET("/filters", func(c *gin.Context) { c.JSON(200, []any{}) })
		v1Api.GET("/featured_tags", func(c *gin.Context) { c.JSON(200, []any{}) })
		v1Api.GET("/scheduled_statuses", func(c *gin.Context) { c.JSON(200, []any{}) })
		v2Api.GET("/filters", func(c *gin.Context) { c.JSON(200, []any{}) })
	}
}
