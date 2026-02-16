package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func MutesRouter(r *gin.RouterGroup) {
	r.GET("/mutes", MutesListHandler)
}

func MutesListHandler(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	var query struct {
		Limit   int    `form:"limit"`
		MaxID   string `form:"max_id"`
		SinceID string `form:"since_id"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	if query.Limit <= 0 {
		query.Limit = 40
	}
	accounts, err := misskey.MutesList(ctx, query.Limit, query.SinceID, query.MaxID)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}
