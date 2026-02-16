package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func AnnouncementsRouter(r *gin.RouterGroup) {
	r.GET("/announcements", AnnouncementsHandler)
	r.POST("/announcements/:id/dismiss", AnnouncementDismissHandler)
}

func AnnouncementsHandler(c *gin.Context) {
	ctx, _ := misstodon.ContextWithGinContext(c)
	announcements, err := misskey.Announcements(ctx, 20)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(announcements))
}

func AnnouncementDismissHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
