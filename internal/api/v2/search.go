package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func SearchRouter(r *gin.RouterGroup) {
	r.GET("/search", SearchHandler)
}

func SearchHandler(c *gin.Context) {
	ctx, _ := misstodon.ContextWithGinContext(c)
	var query struct {
		Q       string `form:"q"`
		Type    string `form:"type"`
		Limit   int    `form:"limit"`
		Offset  int    `form:"offset"`
		Resolve bool   `form:"resolve"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}
	result, err := misskey.Search(ctx, query.Q, query.Type, query.Limit, query.Offset)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
