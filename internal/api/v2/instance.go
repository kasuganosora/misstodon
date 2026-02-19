package v2

import (
	"fmt"
	"net/http"
	"strings"

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
	
	proxyHost := c.Request.Host
	
	// DEBUG: 返回 wxw.moe 的数据，替换域名来测试
	debugResp := `{"domain":"REPLACE_HOST","title":"呜呜 w(> ? <)w","version":"4.3.9~wxw","source_url":"https://github.com/wxwmoe/mastodon","description":"一个萌萌哒 泛ACGN 实例，欢迎安家 ~","usage":{"users":{"active_month":3123}},"thumbnail":{"url":"https://ovo.wxw.moe/site_uploads/files/000/000/004/@1x/ca1c3446c33703ad.png","blurhash":"USR:1rRp?ws%tRV@bJn#%Nt7ROWFxtofa#Rj","versions":{"@1x":"https://ovo.wxw.moe/site_uploads/files/000/000/004/@1x/ca1c3446c33703ad.png","@2x":"https://ovo.wxw.moe/site_uploads/files/000/000/004/@2x/ca1c3446c33703ad.png"}},"icon":[{"src":"https://zzz.wxw.moe/packs/media/icons/android-chrome-36x36-45b0f151edf6ceed682187a0ed02ca59.png","size":"36x36"},{"src":"https://zzz.wxw.moe/packs/media/icons/android-chrome-48x48-1bb30efa6265f7ff942fdd874f37a3ee.png","size":"48x48"},{"src":"https://zzz.wxw.moe/packs/media/icons/android-chrome-72x72-202570e790647a1c9b5b58954ae0a1a5.png","size":"72x72"},{"src":"https://zzz.wxw.moe/packs/media/icons/android-chrome-96x96-323a425dc6f29beaa74eacfe9e000770.png","size":"96x96"},{"src":"https://zzz.wxw.moe/packs/media/icons/android-chrome-144x144-d58a19f58e90e3eb27e1c25eb5bf00a4.png","size":"144x144"},{"src":"https://zzz.wxw.moe/packs/media/icons/android-chrome-192x192-11b9ad6edfb539c2a667994d1d512570.png","size":"192x192"},{"src":"https://zzz.wxw.moe/packs/media/icons/android-chrome-256x256-e49bbed09ec968f1e0a7dd16e6864a4e.png","size":"256x256"},{"src":"https://zzz.wxw.moe/packs/media/icons/android-chrome-384x384-65752f7ee62a3d733d54dff3f7de6d58.png","size":"384x384"},{"src":"https://zzz.wxw.moe/packs/media/icons/android-chrome-512x512-ba16300555530a0a976c9d0a7fcc9dc6.png","size":"512x512"}],"languages":["en"],"configuration":{"urls":{"streaming":"wss://REPLACE_HOST/api/v1/streaming"},"vapid":{"public_key":"BCUeDMIDchElG7FSb9iAq4gtIvCqpJlZv1yZ5QdV0NHy3hBvyw47YA5llwGmdmdBje3sq7vUddyVgJS-y-kL2Kk="},"accounts":{"max_featured_tags":10,"max_pinned_statuses":5},"statuses":{"max_characters":20000,"max_media_attachments":4,"characters_reserved_per_url":23},"media_attachments":{"supported_mime_types":["image/jpeg","image/png","image/gif","image/heic","image/heif","image/webp","image/avif","video/webm","video/mp4","video/quicktime","video/ogg","audio/wave","audio/wav","audio/x-wav","audio/x-pn-wave","audio/vnd.wave","audio/ogg","audio/vorbis","audio/mpeg","audio/mp3","audio/webm","audio/flac","audio/aac","audio/m4a","audio/x-m4a","audio/mp4","audio/3gpp","video/x-ms-asf"],"image_size_limit":103809024,"image_matrix_limit":33177600,"video_size_limit":103809024,"video_frame_rate_limit":120,"video_matrix_limit":8294400},"polls":{"max_options":16,"max_characters_per_option":50,"min_expiration":300,"max_expiration":2629746},"translation":{"enabled":true}},"registrations":{"enabled":true,"approval_required":true,"message":null,"url":null},"api_versions":{"mastodon":2},"contact":{"email":"support@REPLACE_HOST","account":{"id":"3","username":"wxw_moe_status","acct":"wxw_moe_status","display_name":"wxw.moe 更新姬","locked":false,"bot":false,"discoverable":true,"indexable":true,"group":false,"created_at":"2017-11-19T00:00:00.000Z","note":"<p>平时一般不会打扰大家，可以放心关注 ~</p>","url":"https://REPLACE_HOST/@wxw_moe_status","uri":"https://REPLACE_HOST/users/wxw_moe_status","avatar":"https://ovo.wxw.moe/accounts/avatars/000/000/003/original/b66516170bb2dc3e.png","avatar_static":"https://ovo.wxw.moe/accounts/avatars/000/000/003/original/b66516170bb2dc3e.png","header":"https://ovo.wxw.moe/accounts/headers/000/000/003/original/07e1ae74f458e649.jpg","header_static":"https://ovo.wxw.moe/accounts/headers/000/000/003/original/07e1ae74f458e649.jpg","followers_count":12306,"following_count":3,"statuses_count":218,"last_status_at":"2025-12-18","hide_collections":true,"noindex":false,"emojis":[],"roles":[],"fields":[]}},"rules":[]}`
	
	debugResp = strings.ReplaceAll(debugResp, "REPLACE_HOST", proxyHost)
	c.Data(http.StatusOK, "application/json", []byte(debugResp))
}

func InstanceV2HandlerReal(c *gin.Context) {
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
	// Icon should be empty array, not null
	v2.Icon = []models.InstanceIcon{}
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
