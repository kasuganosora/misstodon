package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func TrendsRouter(r *gin.RouterGroup) {
	group := r.Group("/trends")
	group.GET("/tags", TrendsTagsHandler)
	group.GET("/statuses", TrendsStatusHandler)
}

func TrendsTagsHandler(c *gin.Context) {
	limit := 10
	if v, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = v
		if limit > 20 {
			limit = 20
		}
	}
	offset, _ := strconv.Atoi(c.Query("offset"))

	ctx, _ := misstodon.ContextWithGinContext(c)

	tags, err := misskey.TrendsTags(ctx, limit, offset)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(tags))
}

func TrendsStatusHandler(c *gin.Context) {
	limit := 20
	if v, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = v
		if limit > 30 {
			limit = 30
		}
	}
	offset, _ := strconv.Atoi(c.Query("offset"))
	ctx, _ := misstodon.ContextWithGinContext(c)
	statuses, err := misskey.TrendsStatus(ctx, limit, offset)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(statuses))
}
