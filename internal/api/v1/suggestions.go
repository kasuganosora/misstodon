package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func SuggestionsHandler(c *gin.Context) {
	ctx, _ := misstodon.ContextWithGinContext(c)
	var query struct {
		Limit int `form:"limit"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	if query.Limit <= 0 {
		query.Limit = 40
	}
	accounts, err := misskey.Suggestions(ctx, query.Limit)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}
