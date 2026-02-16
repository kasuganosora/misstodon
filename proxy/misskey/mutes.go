package misskey

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func MutesList(ctx Context, limit int, sinceID, maxID string) ([]models.Account, error) {
	var result []struct {
		ID        string        `json:"id"`
		CreatedAt string        `json:"createdAt"`
		MuteeId   string        `json:"muteeId"`
		Mutee     models.MkUser `json:"mutee"`
	}
	body := makeBody(ctx, utils.Map{"limit": limit})
	if sinceID != "" {
		body["sinceId"] = sinceID
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/mute/list"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var accounts []models.Account
	for _, r := range result {
		if a, err := r.Mutee.ToAccount(ctx.ProxyServer()); err == nil {
			accounts = append(accounts, a)
		}
	}
	return accounts, nil
}
