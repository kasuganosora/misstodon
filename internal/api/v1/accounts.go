package v1

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/pkg/errors"
)

func AccountsRouter(r *gin.RouterGroup) {
	group := r.Group("/accounts")
	r.GET("/favourites", AccountFavourites)
	group.GET("/verify_credentials", AccountsVerifyCredentialsHandler)
	group.PATCH("/update_credentials", AccountsUpdateCredentialsHandler)
	group.GET("/search", AccountSearchHandler)
	group.GET("/lookup", AccountsLookupHandler)
	group.GET("/:id", AccountsGetHandler)
	group.GET("/:id/statuses", AccountsStatusesHandler)
	group.GET("/:id/followers", AccountFollowers)
	group.GET("/:id/following", AccountFollowing)
	group.GET("/relationships", AccountRelationships)
	group.POST("/:id/follow", AccountFollow)
	group.POST("/:id/unfollow", AccountUnfollow)
	group.POST("/:id/mute", AccountMute)
	group.POST("/:id/unmute", AccountUnmute)
	group.POST("/:id/block", AccountBlock)
	group.POST("/:id/unblock", AccountUnblock)
	group.GET("/:id/lists", func(c *gin.Context) { c.JSON(200, []any{}) })
	group.GET("/:id/featured_tags", func(c *gin.Context) { c.JSON(200, []any{}) })
}

func AccountsVerifyCredentialsHandler(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	info, err := misskey.VerifyCredentials(ctx)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func AccountsLookupHandler(c *gin.Context) {
	acct := c.Query("acct")
	if acct == "" {
		c.JSON(http.StatusBadRequest, httperror.ServerError{
			Error: "acct is required",
		})
		return
	}
	ctx, _ := misstodon.ContextWithGinContext(c)
	info, err := misskey.AccountsLookup(ctx, acct)
	if err != nil {
		if errors.Is(err, misskey.ErrNotFound) {
			c.JSON(http.StatusNotFound, httperror.ServerError{
				Error: "Record not found",
			})
			return
		} else if errors.Is(err, misskey.ErrAcctIsInvalid) {
			c.JSON(http.StatusBadRequest, httperror.ServerError{
				Error: err.Error(),
			})
			return
		}
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	if info.Header == "" || info.HeaderStatic == "" {
		info.Header = fmt.Sprintf("https://%s/static/missing.png", c.Request.Host)
		info.HeaderStatic = info.Header
	}
	c.JSON(http.StatusOK, info)
}

func AccountsStatusesHandler(c *gin.Context) {
	uid := c.Param("id")

	ctx, _ := misstodon.ContextWithGinContext(c)

	limit := 30
	pinnedOnly := false
	onlyMedia := false
	onlyPublic := false
	excludeReplies := false
	excludeReblogs := false
	maxID := ""
	minID := ""

	var query struct {
		Limit         *int  `form:"limit"`
		PinnedOnly    *bool `form:"pinned_only"`
		OnlyMedia     *bool `form:"only_media"`
		OnlyPublic    *bool `form:"only_public"`
		ExcludeReplies *bool `form:"exclude_replies"`
		ExcludeReblogs *bool `form:"exclude_reblogs"`
		MaxID         string `form:"max_id"`
		MinID         string `form:"min_id"`
	}
	if err := c.ShouldBindQuery(&query); err == nil {
		if query.Limit != nil {
			limit = *query.Limit
		}
		if query.PinnedOnly != nil {
			pinnedOnly = *query.PinnedOnly
		}
		if query.OnlyMedia != nil {
			onlyMedia = *query.OnlyMedia
		}
		if query.OnlyPublic != nil {
			onlyPublic = *query.OnlyPublic
		}
		if query.ExcludeReplies != nil {
			excludeReplies = *query.ExcludeReplies
		}
		if query.ExcludeReblogs != nil {
			excludeReblogs = *query.ExcludeReblogs
		}
		maxID = query.MaxID
		minID = query.MinID
	}

	statuses, err := misskey.AccountsStatuses(
		ctx, uid,
		limit,
		pinnedOnly, onlyMedia, onlyPublic, excludeReplies, excludeReblogs,
		maxID, minID)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(statuses))
}

func AccountsUpdateCredentialsHandler(c *gin.Context) {
	form, err := parseAccountsUpdateCredentialsForm(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}

	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}

	account, err := misskey.UpdateCredentials(ctx,
		form.DisplayName, form.Note,
		form.Locked, form.Bot, form.Discoverable,
		form.SourcePrivacy, form.SourceSensitive, form.SourceLanguage,
		form.AccountFields,
		form.Avatar, form.Header)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, account)
}

type accountsUpdateCredentialsForm struct {
	DisplayName     *string `form:"display_name"`
	Note            *string `form:"note"`
	Locked          *bool   `form:"locked"`
	Bot             *bool   `form:"bot"`
	Discoverable    *bool   `form:"discoverable"`
	SourcePrivacy   *string `form:"source[privacy]"`
	SourceSensitive *bool   `form:"source[sensitive]"`
	SourceLanguage  *string `form:"source[language]"`
	AccountFields   []models.AccountField
	Avatar          *multipart.FileHeader
	Header          *multipart.FileHeader
}

func parseAccountsUpdateCredentialsForm(c *gin.Context) (f accountsUpdateCredentialsForm, err error) {
	var form accountsUpdateCredentialsForm
	if err = c.ShouldBind(&form); err != nil {
		return
	}

	var values = make(map[string][]string)
	for k, v := range c.Request.URL.Query() {
		values[k] = v
	}
	if c.Request.Method == "POST" || c.Request.Method == "PATCH" {
		if fp := c.Request.PostForm; fp != nil {
			for k, v := range fp {
				values[k] = v
			}
		}
	}
	if mf := c.Request.MultipartForm; mf != nil {
		for k, v := range mf.Value {
			values[k] = v
		}
	}
	for _, field := range utils.GetFieldsAttributes(values) {
		form.AccountFields = append(form.AccountFields, models.AccountField(field))
	}
	if fh, err := c.FormFile("avatar"); err == nil {
		form.Avatar = fh
	}
	if fh, err := c.FormFile("header"); err == nil {
		form.Header = fh
	}
	return form, nil
}

func AccountFollowRequests(c *gin.Context) {
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
	if err = c.ShouldBindQuery(&query); err != nil {
		httperror.AbortWithError(c, http.StatusBadRequest, err)
		return
	}
	if query.Limit <= 0 {
		query.Limit = 40
	}
	accounts, err := misskey.AccountFollowRequests(ctx, query.Limit, query.SinceID, query.MaxID)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}

func FollowRequestAuthorize(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	id := c.Param("id")
	if err := misskey.AccountFollowRequestsAccept(ctx, id); err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{id})
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, relationships[0])
}

func FollowRequestReject(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	id := c.Param("id")
	if err := misskey.AccountFollowRequestsReject(ctx, id); err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{id})
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, relationships[0])
}

func AccountFollowers(c *gin.Context) {
	ctx, _ := misstodon.ContextWithGinContext(c)
	id := c.Param("id")
	var query struct {
		Limit   int    `form:"limit"`
		MaxID   string `form:"max_id"`
		MinID   string `form:"min_id"`
		SinceID string `form:"since_id"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	if query.Limit <= 0 {
		query.Limit = 40
	}
	if query.Limit > 80 {
		query.Limit = 80
	}
	accounts, err := misskey.AccountFollowers(ctx, id, query.Limit, query.SinceID, query.MinID, query.MaxID)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}

func AccountFollowing(c *gin.Context) {
	ctx, _ := misstodon.ContextWithGinContext(c)

	id := c.Param("id")
	var query struct {
		Limit   int    `form:"limit"`
		MaxID   string `form:"max_id"`
		MinID   string `form:"min_id"`
		SinceID string `form:"since_id"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	if query.Limit <= 0 {
		query.Limit = 40
	}
	if query.Limit > 80 {
		query.Limit = 80
	}
	accounts, err := misskey.AccountFollowing(ctx, id, query.Limit, query.SinceID, query.MinID, query.MaxID)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}

func AccountRelationships(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	var ids []string
	for k, v := range c.Request.URL.Query() {
		if k == "id[]" {
			ids = append(ids, v...)
			continue
		}
	}
	relationships, err := misskey.AccountRelationships(ctx, ids)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, relationships)
}

func AccountFollow(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	id := c.Param("id")
	if err = misskey.AccountFollow(ctx, id); err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{id})
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, relationships[0])
}

func AccountUnfollow(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	id := c.Param("id")
	if err = misskey.AccountUnfollow(ctx, id); err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{id})
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, relationships[0])
}

func AccountMute(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	var params struct {
		ID       string `uri:"id"`
		Duration int64  `json:"duration" form:"duration"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		httperror.AbortWithError(c, http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&params); err != nil {
		httperror.AbortWithError(c, http.StatusBadRequest, err)
		return
	}
	if err = misskey.AccountMute(ctx, params.ID, params.Duration); err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{params.ID})
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, relationships[0])
}

func AccountUnmute(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	id := c.Param("id")
	if err = misskey.AccountUnmute(ctx, id); err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{id})
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, relationships[0])
}

func AccountsGetHandler(c *gin.Context) {
	ctx, _ := misstodon.ContextWithGinContext(c)
	info, err := misskey.AccountsGet(ctx, c.Param("id"))
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	if info.Header == "" || info.HeaderStatic == "" {
		info.Header = fmt.Sprintf("https://%s/static/missing.png", c.Request.Host)
		info.HeaderStatic = info.Header
	}
	c.JSON(http.StatusOK, info)
}

func AccountFavourites(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}

	var params struct {
		Limit   int    `form:"limit"`
		MaxID   string `form:"max_id"`
		MinID   string `form:"min_id"`
		SinceID string `form:"since_id"`
	}
	if err = c.ShouldBindQuery(&params); err != nil {
		httperror.AbortWithError(c, http.StatusBadRequest, err)
		return
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	list, err := misskey.AccountFavourites(ctx,
		params.Limit, params.SinceID, params.MinID, params.MaxID)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(list))
}

func AccountSearchHandler(c *gin.Context) {
	ctx, _ := misstodon.ContextWithGinContext(c)
	var query struct {
		Q       string `form:"q"`
		Limit   int    `form:"limit"`
		Offset  int    `form:"offset"`
		Resolve bool   `form:"resolve"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httperror.ServerError{Error: err.Error()})
		return
	}
	if query.Limit <= 0 {
		query.Limit = 40
	}
	accounts, err := misskey.AccountSearch(ctx, query.Q, query.Limit, query.Offset)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, utils.SliceIfNull(accounts))
}

func AccountBlock(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	id := c.Param("id")
	if err = misskey.AccountBlock(ctx, id); err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{id})
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, relationships[0])
}

func AccountUnblock(c *gin.Context) {
	ctx, err := misstodon.ContextWithGinContext(c, true)
	if err != nil {
		httperror.AbortWithError(c, http.StatusUnauthorized, err)
		return
	}
	id := c.Param("id")
	if err = misskey.AccountUnblock(ctx, id); err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	relationships, err := misskey.AccountRelationships(ctx, []string{id})
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, relationships[0])
}
