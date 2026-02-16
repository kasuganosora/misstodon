package v1

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/pkg/errors"
)

func ApplicationRouter(r *gin.RouterGroup) {
	group := r.Group("/apps")
	group.POST("", ApplicationCreateHandler)
	group.GET("/verify_credentials", ApplicationVerifyCredentials)
}

func ApplicationVerifyCredentials(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":      "Misstodon",
		"vapid_key": "",
	})
}

func ApplicationCreateHandler(c *gin.Context) {
	var params struct {
		ClientName   string `json:"client_name" form:"client_name"`
		WebSite      string `json:"website" form:"website"`
		RedirectUris string `json:"redirect_uris" form:"redirect_uris"`
		Scopes       string `json:"scopes" form:"scopes"`
	}
	if err := c.ShouldBind(&params); err != nil {
		httperror.AbortWithError(c, http.StatusBadRequest, err)
		return
	}
	if params.ClientName == "" || params.RedirectUris == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_name and redirect_uris are required"})
		return
	}
	server := c.GetString("proxy-server")
	u, err := url.Parse(strings.Join([]string{"https://", c.Request.Host, "/oauth/redirect"}, ""))
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, errors.WithStack(err))
		return
	}
	query := u.Query()
	query.Add("server", server)
	query.Add("redirect_uris", params.RedirectUris)
	u.RawQuery = query.Encode()
	app, err := misskey.ApplicationCreate(
		server,
		params.ClientName,
		u.String(),
		params.Scopes,
		params.WebSite)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	// Encode the secret into client_id so we can retrieve it in the authorize step
	// without needing a database. Format: "realId.secret"
	if app.ClientID != nil && app.ClientSecret != nil {
		encodedID := *app.ClientID + "." + *app.ClientSecret
		app.ClientID = &encodedID
	}
	c.JSON(http.StatusOK, app)
}
