package v2

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/internal/api/httperror"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
)

func InstanceRouter(r *gin.RouterGroup) {
	r.GET("/instance", InstanceV2Handler)
}

func InstanceV2Handler(c *gin.Context) {
	server := c.GetString("proxy-server")
	fmt.Printf("[DEBUG] InstanceV2Handler called, server=%s\n", server)
	info, err := misskey.Instance(server, global.AppVersion, c.Request.Host)
	if err != nil {
		fmt.Printf("[DEBUG] Instance error: %v\n", err)
		httperror.AbortWithError(c, http.StatusInternalServerError, err)
		return
	}

	proxyHost := c.Request.Host
	v2 := models.InstanceV2{
		Domain:      info.Uri,
		Title:       info.Title,
		Version:     info.Version,
		SourceURL:   "https://github.com/gizmo-ds/misstodon",
		Description: info.Description,
	}
	v2.Usage.Users.ActiveMonth = info.Stats.UserCount
	v2.Thumbnail.URL = info.Thumbnail
	// Streaming URL - important for client compatibility
	v2.Configuration.Urls.Streaming = "wss://" + proxyHost + "/api/v1/streaming"
	if langs, ok := info.Languages.([]string); ok {
		v2.Languages = langs
	} else {
		v2.Languages = []string{}
	}
	// Configuration
	v2.Configuration.Accounts.MaxFeaturedTags = 10
	v2.Configuration.Accounts.MaxPinnedStatuses = 5
	v2.Configuration.Statuses.MaxCharacters = info.Configuration.Statuses.MaxCharacters
	v2.Configuration.Statuses.MaxMediaAttachments = info.Configuration.Statuses.MaxMediaAttachments
	v2.Configuration.Statuses.CharactersReservedPerUrl = info.Configuration.Statuses.CharactersReservedPerUrl
	v2.Configuration.MediaAttachments.SupportedMimeTypes = info.Configuration.MediaAttachments.SupportedMimeTypes
	v2.Configuration.MediaAttachments.ImageSizeLimit = 10485760
	v2.Configuration.MediaAttachments.ImageMatrixLimit = 16777216
	v2.Configuration.MediaAttachments.VideoSizeLimit = 41943040
	v2.Configuration.MediaAttachments.VideoFrameRateLimit = 60
	v2.Configuration.MediaAttachments.VideoMatrixLimit = 2304000
	v2.Configuration.Polls.MaxOptions = 10
	v2.Configuration.Polls.MaxCharactersPerOption = 50
	v2.Configuration.Polls.MinExpiration = 300
	v2.Configuration.Polls.MaxExpiration = 2629746
	v2.Configuration.Translation.Enabled = false
	// Registrations
	v2.Registrations.Enabled = info.Registrations
	v2.Registrations.ApprovalRequired = false
	// API Versions (important for client compatibility)
	v2.ApiVersions.Mastodon = 2
	// Contact
	v2.Contact.Email = info.Email
	v2.Rules = info.Rules
	if v2.Rules == nil {
		v2.Rules = []models.InstanceRule{}
	}
	c.JSON(http.StatusOK, v2)
}
