package misskey

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// SupportedMimeTypes is a list of supported mime types
//
// https://github.com/misskey-dev/misskey/blob/79212bbd375705f0fd658dd5b50b47f77d622fb8/packages/backend/src/const.ts#L25
var SupportedMimeTypes = []string{
	"image/png",
	"image/gif",
	"image/jpeg",
	"image/webp",
	"image/avif",
	"image/apng",
	"image/bmp",
	"image/tiff",
	"image/x-icon",
	"audio/opus",
	"video/ogg",
	"audio/ogg",
	"application/ogg",
	"video/quicktime",
	"video/mp4",
	"audio/mp4",
	"video/x-m4v",
	"audio/x-m4a",
	"video/3gpp",
	"video/3gpp2",
	"video/mpeg",
	"audio/mpeg",
	"video/webm",
	"audio/webm",
	"audio/aac",
	"audio/x-flac",
	"audio/vnd.wave",
}

func Instance(server, version, proxyHost string) (models.Instance, error) {
	var info models.Instance
	var serverInfo models.MkMeta
	apiURL := utils.JoinURL(server, "/api/meta")
	fmt.Printf("[DEBUG] Instance: server=%s, apiURL=%s\n", server, apiURL)
	resp, err := client.R().
		SetBody(map[string]any{
			"detail": false,
		}).
		SetResult(&serverInfo).
		Post(apiURL)
	if err != nil {
		fmt.Printf("[DEBUG] /api/meta error: %v\n", err)
		log.Error().Err(err).Str("server", server).Str("url", apiURL).Msg("Failed to call /api/meta")
		return info, err
	}
	fmt.Printf("[DEBUG] /api/meta status: %d, body: %s\n", resp.StatusCode(), string(resp.Body()))
	serverUrl, err := url.Parse(serverInfo.URI)
	if err != nil {
		log.Error().Err(err).Str("uri", serverInfo.URI).Msg("Failed to parse server URI")
		return info, err
	}
	if resp.StatusCode() != http.StatusOK {
		log.Error().Int("status", resp.StatusCode()).Str("body", string(resp.Body())).Msg("Non-OK response from /api/meta")
		return info, errors.New("Failed to get instance info")
	}
	domain := serverUrl.Host
	if proxyHost != "" {
		domain = proxyHost
	}
	info = models.Instance{
		Uri:              domain,
		Title:            serverInfo.Name,
		Description:      serverInfo.Description,
		ShortDescription: serverInfo.Description,
		Email:            serverInfo.MaintainerEmail,
		Version:          version,
		Thumbnail:        serverInfo.BannerUrl,
		Registrations:    !serverInfo.DisableRegistration,
		InvitesEnabled:   serverInfo.Policies.CanInvite,
		Rules:            []models.InstanceRule{},
		Languages:        serverInfo.Langs,
	}
	// TODO: 需要先实现 `/streaming`
	// info.Urls.StreamingApi = serverInfo.StreamingAPI
	if info.Languages == nil {
		info.Languages = []string{}
	}
	info.Configuration.Statuses.MaxCharacters = serverInfo.MaxNoteTextLength
	// NOTE: misskey没有相关限制, 此处返回固定值
	info.Configuration.Statuses.MaxMediaAttachments = 4
	// NOTE: misskey没有相关设置, 此处返回固定值
	info.Configuration.Statuses.CharactersReservedPerUrl = 23
	info.Configuration.Accounts.MaxFeaturedTags = 10
	info.Configuration.MediaAttachments.SupportedMimeTypes = SupportedMimeTypes
	info.Configuration.MediaAttachments.ImageSizeLimit = 10485760
	info.Configuration.MediaAttachments.ImageMatrixLimit = 16777216
	info.Configuration.MediaAttachments.VideoSizeLimit = 41943040
	info.Configuration.MediaAttachments.VideoFrameRateLimit = 60
	info.Configuration.MediaAttachments.VideoMatrixLimit = 2304000
	info.Configuration.Polls.MaxOptions = 10
	info.Configuration.Polls.MaxCharactersPerOption = 50
	info.Configuration.Polls.MinExpiration = 300
	info.Configuration.Polls.MaxExpiration = 2629746

	var serverStats models.MkStats
	statsURL := utils.JoinURL(server, "/api/stats")
	resp, err = client.R().
		SetBody(map[string]any{}).
		SetResult(&serverStats).
		Post(statsURL)
	if err != nil {
		log.Error().Err(err).Str("url", statsURL).Msg("Failed to call /api/stats")
		return info, err
	}
	if resp.StatusCode() != http.StatusOK {
		log.Error().Int("status", resp.StatusCode()).Str("body", string(resp.Body())).Msg("Non-OK response from /api/stats")
		return info, errors.New("Failed to get instance info")
	}
	info.Stats.UserCount = serverStats.OriginalUsersCount
	info.Stats.StatusCount = serverStats.OriginalNotesCount
	info.Stats.DomainCount = serverStats.Instances
	return info, err
}

func InstancePeers(server string) ([]string, error) {
	return nil, nil
}

func InstanceCustomEmojis(server string) ([]models.CustomEmoji, error) {
	var emojis struct {
		Emojis []models.MkEmoji `json:"emojis"`
	}
	resp, err := client.R().
		SetResult(&emojis).
		SetBody(utils.Map{}).
		Post(utils.JoinURL(server, "/api/emojis"))
	if err != nil {
		return nil, err
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	return lo.Map(emojis.Emojis, func(e models.MkEmoji, _ int) models.CustomEmoji {
		return e.ToCustomEmoji()
	}), nil
}
