package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func TimelinesRouter(r *gin.RouterGroup) {
	group := r.Group("/timelines")
	group.GET("/public", TimelinePublicHandler)
	group.GET("/home", TimelineHomeHandler)
	group.GET("/tag/:hashtag", TimelineHashtag)
}

func TimelinePublicHandler(c *gin.Context) {
	ctx, _ := misstodon.ContextWithGinContext(c)
	limit := 20
	if v, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = v
		if limit > 40 {
			limit = 40
		}
	}
	timelineType := models.TimelinePublicTypeRemote
	if c.Query("local") == "true" {
		timelineType = models.TimelinePublicTypeLocal
	}
	list, err := misskey.TimelinePublic(ctx,
		timelineType, c.Query("only_media") == "true", limit,
		c.Query("max_id"), c.Query("min_id"))
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(list))
}

func TimelineHomeHandler(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	limit := 20
	if v, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = v
		if limit > 40 {
			limit = 40
		}
	}
	list, err := misskey.TimelineHome(ctx,
		limit, c.Query("max_id"), c.Query("min_id"))
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(list))
}

func TimelineHashtag(c *gin.Context) {
	ctx, _ := misstodon.ContextWithGinContext(c)

	limit := 20
	if v, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = v
		if limit > 40 {
			limit = 40
		}
	}

	list, err := misskey.SearchStatusByHashtag(ctx, c.Param("hashtag"),
		limit, c.Query("max_id"), c.Query("since_id"), c.Query("min_id"))
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(list))
}
