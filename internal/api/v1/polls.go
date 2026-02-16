package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func PollsRouter(r *gin.RouterGroup) {
	group := r.Group("/polls")
	group.GET("/:id", PollGetHandler)
	group.POST("/:id/votes", PollVoteHandler)
}

func PollGetHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, _ := misstodon.ContextWithGinContext(c)
	status, err := misskey.StatusSingle(ctx, id)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	if status.Poll == nil {
		c.JSON(http.StatusNotFound, httperror.ServerError{Error: "Record not found"})
		return
	}
	c.JSON(http.StatusOK, status.Poll)
}

func PollVoteHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	var form struct {
		Choices []int `json:"choices"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	poll, err := misskey.PollVote(ctx, id, form.Choices)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, poll)
}
