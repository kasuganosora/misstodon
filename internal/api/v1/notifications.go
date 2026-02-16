package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/samber/lo"
)

func NotificationsRouter(r *gin.RouterGroup) {
	group := r.Group("/notifications")
	group.GET("", NotificationsHandler)
	group.GET("/unread_count", NotificationsUnreadCount)
	group.POST("/clear", NotificationsClear)
	group.GET("/:id", NotificationGet)
	group.POST("/:id/dismiss", NotificationDismiss)
}

func NotificationsUnreadCount(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		// Fall back to 0 if not authenticated
		c.JSON(http.StatusOK, gin.H{"count": 0})
		return
	}
	count, err := misskey.NotificationsUnreadCount(ctx)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"count": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func NotificationsClear(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	if err := misskey.NotificationsClear(ctx); err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func NotificationGet(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	notification, err := misskey.NotificationGet(ctx, id)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, notification)
}

func NotificationDismiss(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	_ = ctx
	c.JSON(http.StatusOK, gin.H{})
}

func NotificationsHandler(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	var query struct {
		MaxId   string `form:"max_id"`
		MinId   string `form:"min_id"`
		SinceId string `form:"since_id"`
		Limit   int    `form:"limit"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}

	getTypes := func(name string) []models.NotificationType {
		types := lo.Map(c.QueryArray(name), func(item string, _ int) models.NotificationType { return models.NotificationType(item) })
		types = lo.Filter(types, func(item models.NotificationType, _ int) bool {
			return item != "" && item.ToMkNotificationType() != models.MkNotificationTypeUnknown
		})
		return types
	}

	types := getTypes("types[]")
	excludeTypes := getTypes("exclude_types[]")

	result, err := misskey.NotificationsGet(ctx,
		query.Limit, query.SinceId, query.MinId, query.MaxId,
		types, excludeTypes, "")
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
