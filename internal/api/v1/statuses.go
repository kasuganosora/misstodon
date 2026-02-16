package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/pkg/errors"
)

func StatusesRouter(r *gin.RouterGroup) {
	group := r.Group("/statuses")
	group.POST("", PostNewStatus)
	group.GET("/:id", StatusHandler)
	group.DELETE("/:id", StatusDeleteHandler)
	group.GET("/:id/context", StatusContextHandler)
	group.GET("/:id/reblogged_by", StatusRebloggedByHandler)
	group.GET("/:id/favourited_by", StatusFavouritedByHandler)
	group.POST("/:id/bookmark", StatusBookmark)
	group.POST("/:id/unbookmark", StatusUnBookmark)
	group.POST("/:id/favourite", StatusFavourite)
	group.POST("/:id/unfavourite", StatusUnFavourite)
	group.POST("/:id/reblog", StatusReblogHandler)
	group.POST("/:id/unreblog", StatusUnreblogHandler)
	group.POST("/:id/mute", StatusMuteHandler)
	group.POST("/:id/unmute", StatusUnmuteHandler)
	group.POST("/:id/pin", StatusPinHandler)
	group.POST("/:id/unpin", StatusUnpinHandler)
	group.POST("/:id/translate", StatusTranslateHandler)
}

func StatusHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, _ := misstodon.ContextWithGinContext(c)
	info, err := misskey.StatusSingle(ctx, id)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func StatusContextHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, _ := misstodon.ContextWithGinContext(c)
	context, err := misskey.StatusContext(ctx, id)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, context)
}

func StatusDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusSingle(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrNotFound) {
			c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
			return
		}
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	if err := misskey.StatusDelete(ctx, id); err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func StatusReblogHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusReblog(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		}
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func StatusUnreblogHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusUnreblog(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		}
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func StatusRebloggedByHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, _ := misstodon.ContextWithGinContext(c)
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
	accounts, err := misskey.StatusRebloggedBy(ctx, id, query.Limit, query.SinceID, query.MaxID)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}

func StatusFavouritedByHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, _ := misstodon.ContextWithGinContext(c)
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
	accounts, err := misskey.StatusFavouritedBy(ctx, id, query.Limit, query.SinceID, query.MaxID)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}

func StatusMuteHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusMuteThread(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		}
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func StatusUnmuteHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusUnmuteThread(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		}
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func StatusPinHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusPin(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		}
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func StatusUnpinHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusUnpin(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		}
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func StatusTranslateHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	var form struct {
		Lang string `json:"lang" form:"lang"`
	}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	if form.Lang == "" {
		form.Lang = "en"
	}
	translation, err := misskey.StatusTranslate(ctx, id, form.Lang)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, translation)
}

func StatusFavourite(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusFavourite(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		} else if errors.Is(err, misskey.ErrNotFound) {
			c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
			return
		} else {
			httperror.AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
	}
	c.JSON(http.StatusOK, status)
}

func StatusUnFavourite(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusUnFavourite(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		} else if errors.Is(err, misskey.ErrNotFound) {
			c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
			return
		} else {
			httperror.AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
	}
	c.JSON(http.StatusOK, status)
}

func StatusBookmark(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusBookmark(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		} else if errors.Is(err, misskey.ErrNotFound) {
			c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
			return
		} else {
			httperror.AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
	}
	c.JSON(http.StatusOK, status)
}

func StatusUnBookmark(c *gin.Context) {
	id := c.Param("id")
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	status, err := misskey.StatusUnBookmark(ctx, id)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		} else if errors.Is(err, misskey.ErrNotFound) {
			c.JSON(http.StatusNotFound, httperror.ServerError{Error: err.Error()})
			return
		} else {
			httperror.AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
	}
	c.JSON(http.StatusOK, status)
}

func StatusBookmarks(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	var query struct {
		Limit   int    `form:"limit"`
		MaxID   string `form:"max_id"`
		MinID   string `form:"min_id"`
		SinceID string `form:"since_id"`
	}
	if err = c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}
	status, err := misskey.StatusBookmarks(ctx, query.Limit, query.SinceID, query.MinID, query.MaxID)
	if err != nil {
		if errors.Is(err, misskey.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, httperror.ServerError{Error: err.Error()})
			return
		} else {
			httperror.AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(status))
}

type postNewStatusPollForm struct {
	Options    []string `json:"options"`
	ExpiresIn  int      `json:"expires_in"`
	Multiple   bool     `json:"multiple"`
	HideTotals bool     `json:"hide_totals"`
}

type postNewStatusForm struct {
	Status      *string                 `json:"status"`
	Poll        *postNewStatusPollForm  `json:"poll"`
	MediaIDs    []string                `json:"media_ids"`
	InReplyToID string                  `json:"in_reply_to_id"`
	Sensitive   bool                    `json:"sensitive"`
	SpoilerText string                  `json:"spoiler_text"`
	Visibility  models.StatusVisibility `json:"visibility"`
	Language    string                  `json:"language"`
	ScheduledAt time.Time               `json:"scheduled_at"`
}

func PostNewStatus(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}

	var form postNewStatusForm
	if err = c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	var pollOptions []string
	var pollExpiresIn int
	var pollMultiple bool
	if form.Poll != nil {
		pollOptions = form.Poll.Options
		pollExpiresIn = form.Poll.ExpiresIn
		pollMultiple = form.Poll.Multiple
	}
	status, err := misskey.PostNewStatus(ctx,
		form.Status, pollOptions, pollExpiresIn, pollMultiple,
		form.MediaIDs, form.InReplyToID,
		form.Sensitive, form.SpoilerText,
		form.Visibility, form.Language,
		form.ScheduledAt)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, status)
}
