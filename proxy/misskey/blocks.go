package misskey

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func AccountBlock(ctx Context, userID string) error {
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"userId": userID})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/blocking/create"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK, "ALREADY_BLOCKING"); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AccountUnblock(ctx Context, userID string) error {
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"userId": userID})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/blocking/delete"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK, "NOT_BLOCKING"); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func BlocksList(ctx Context, limit int, sinceID, maxID string) ([]models.Account, error) {
	var result []struct {
		ID        string        `json:"id"`
		CreatedAt string        `json:"createdAt"`
		BlockeeId string        `json:"blockeeId"`
		Blockee   models.MkUser `json:"blockee"`
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
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/blocking/list"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var accounts []models.Account
	for _, r := range result {
		if a, err := r.Blockee.ToAccount(ctx.ProxyServer()); err == nil {
			accounts = append(accounts, a)
		}
	}
	return accounts, nil
}
