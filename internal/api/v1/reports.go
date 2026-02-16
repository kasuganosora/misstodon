package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func ReportsRouter(r *gin.RouterGroup) {
	r.POST("/reports", ReportCreateHandler)
}

func ReportCreateHandler(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	var form struct {
		AccountID string   `json:"account_id"`
		StatusIDs []string `json:"status_ids"`
		Comment   string   `json:"comment"`
		Forward   bool     `json:"forward"`
		Category  string   `json:"category"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	report, err := misskey.ReportCreate(ctx, form.AccountID, form.Comment, form.StatusIDs)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, report)
}
