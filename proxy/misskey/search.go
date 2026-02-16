package misskey

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func AccountSearch(ctx Context, q string, limit, offset int) ([]models.Account, error) {
	var result []models.MkUser
	body := makeBody(ctx, utils.Map{"query": q, "limit": limit, "offset": offset})
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/users/search"))
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

func SearchStatuses(ctx Context, q string, limit, offset int) ([]models.Status, error) {
	var result []models.MkNote
	body := makeBody(ctx, utils.Map{"query": q, "limit": limit, "offset": offset})
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/search"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	return lo.Map(result, func(n models.MkNote, _ int) models.Status {
		return n.ToStatus(ctx.ProxyServer())
	}), nil
}

func SearchHashtags(ctx Context, q string, limit, offset int) ([]models.Tag, error) {
	var result []string
	body := makeBody(ctx, utils.Map{"query": q, "limit": limit, "offset": offset})
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/hashtags/search"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	return lo.Map(result, func(name string, _ int) models.Tag {
		return models.Tag{
			Name: name,
			Url:  utils.JoinURL(ctx.ProxyServer(), "/tags/", name),
		}
	}), nil
}

func Search(ctx Context, q, searchType string, limit, offset int) (models.SearchResult, error) {
	result := models.SearchResult{
		Accounts: []models.Account{},
		Statuses: []models.Status{},
		Hashtags: []models.Tag{},
	}
	if searchType == "" || searchType == "accounts" {
		if accounts, err := AccountSearch(ctx, q, limit, offset); err == nil && accounts != nil {
			result.Accounts = accounts
		}
	}
	if searchType == "" || searchType == "statuses" {
		if statuses, err := SearchStatuses(ctx, q, limit, offset); err == nil && statuses != nil {
			result.Statuses = statuses
		}
	}
	if searchType == "" || searchType == "hashtags" {
		if tags, err := SearchHashtags(ctx, q, limit, offset); err == nil && tags != nil {
			result.Hashtags = tags
		}
	}
	return result, nil
}
