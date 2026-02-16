package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func MediaRouter(r *gin.RouterGroup) {
	group := r.Group("/media")
	group.POST("", MediaUploadHandler)
}

func MediaUploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	description := c.PostForm("description")

	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}

	ma, err := misskey.MediaUpload(ctx, file, description)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ma)
}
