package misskey

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func Suggestions(ctx Context, limit int) ([]models.Account, error) {
	var result []models.MkUser
	body := makeBody(ctx, utils.Map{"limit": limit})
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/users/recommendation"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var accounts []models.Account
	for _, u := range result {
		if a, err := u.ToAccount(ctx.ProxyServer()); err == nil {
			accounts = append(accounts, a)
		}
	}
	return accounts, nil
}
