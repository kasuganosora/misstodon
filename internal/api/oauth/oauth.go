package oauth

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/api/middleware"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func Router(r *gin.RouterGroup) {
	group := r.Group("/oauth")
	group.Use(middleware.CORS())
	group.GET("/authorize", AuthorizeHandler)
	group.POST("/token", TokenHandler)
	// NOTE: This is not a standard endpoint
	group.GET("/redirect", RedirectHandler)
}

func RedirectHandler(c *gin.Context) {
	redirectUris := c.Query("redirect_uris")
	server := c.Query("server")
	token := c.Query("token")
	if redirectUris == "" || server == "" {
		c.String(http.StatusBadRequest, "redirect_uris and server are required")
		return
	}
	if token == "" {
		if strings.Contains(redirectUris, "?token=") {
			i := strings.Index(redirectUris, "?token=")
			token = redirectUris[i+7:]
			redirectUris = redirectUris[:i]
		}
		if strings.Contains(server, "?token=") {
			i := strings.Index(server, "?token=")
			token = server[i+7:]
			server = server[:i]
		}
	}
	u, err := url.Parse(redirectUris)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	query := u.Query()
	query.Add("code", token)
	u.RawQuery = query.Encode()
	c.Redirect(http.StatusFound, u.String())
}

func TokenHandler(c *gin.Context) {
	var params struct {
		GrantType    string `json:"grant_type" form:"grant_type"`
		ClientID     string `json:"client_id" form:"client_id"`
		ClientSecret string `json:"client_secret" form:"client_secret"`
		RedirectURI  string `json:"redirect_uri" form:"redirect_uri"`
		Code         string `json:"code" form:"code"`
		Scope        string `json:"scope" form:"scope"`
	}
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if params.GrantType == "" || params.ClientID == "" ||
		params.ClientSecret == "" || params.RedirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "grant_type, client_id, client_secret and redirect_uri are required",
		})
		return
	}
	server := c.GetString("proxy-server")
	accessToken, userID, err := misskey.OAuthToken(server, params.Code, params.ClientSecret)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": strings.Join([]string{userID, accessToken}, "."),
		"token_type":   "Bearer",
		"scope":        params.Scope,
		"created_at":   time.Now().Unix(),
	})
}

func AuthorizeHandler(c *gin.Context) {
	var params struct {
		ClientID     string `form:"client_id"`
		RedirectUri  string `form:"redirect_uri"`
		ResponseType string `form:"response_type"`
		Scope        string `form:"scope"`
		Lang         string `form:"lang"`
		ForceLogin   bool   `form:"force_login"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if params.ResponseType != "code" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "response_type must be code",
		})
		return
	}
	if params.ClientID == "" || params.RedirectUri == "" || params.ResponseType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client_id, redirect_uri and response_type are required",
		})
		return
	}
	// Extract the secret from the encoded client_id (format: "realId.secret")
	dotIndex := strings.Index(params.ClientID, ".")
	if dotIndex < 0 || dotIndex >= len(params.ClientID)-1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client_id is invalid",
		})
		return
	}
	secret := params.ClientID[dotIndex+1:]
	server := c.GetString("proxy-server")
	u, err := misskey.OAuthAuthorize(server, secret)
	if err != nil {
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, u)
}
