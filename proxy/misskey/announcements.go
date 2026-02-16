package misskey

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/mfm"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func Announcements(ctx Context, limit int) ([]models.Announcement, error) {
	type mkAnnouncement struct {
		ID        string  `json:"id"`
		CreatedAt string  `json:"createdAt"`
		UpdatedAt *string `json:"updatedAt"`
		Text      string  `json:"text"`
		Title     string  `json:"title"`
		IsRead    bool    `json:"isRead"`
	}
	var result []mkAnnouncement
	body := makeBody(ctx, utils.Map{"limit": limit})
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/announcements"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var announcements []models.Announcement
	for _, a := range result {
		content := a.Text
		if html, err := mfm.ToHtml(a.Text, mfm.Option{
			Url: utils.JoinURL(ctx.ProxyServer()),
		}); err == nil {
			content = html
		}
		ann := models.Announcement{
			ID:          a.ID,
			Content:     content,
			Published:   true,
			AllDay:      false,
			PublishedAt: a.CreatedAt,
			UpdatedAt:   a.CreatedAt,
			Read:        a.IsRead,
			Mentions:    []struct{}{},
			Statuses:    []struct{}{},
			Tags:        []models.Tag{},
			Emojis:      []struct{}{},
			Reactions:   []struct{}{},
		}
		if a.UpdatedAt != nil {
			ann.UpdatedAt = *a.UpdatedAt
		}
		announcements = append(announcements, ann)
	}
	return announcements, nil
}
