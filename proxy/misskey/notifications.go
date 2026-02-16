package misskey

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func NotificationsGet(ctx Context,
	limit int, sinceId, minId, maxId string,
	types, excludeTypes []models.NotificationType, accountId string,
) ([]models.Notification, error) {
	limit = utils.NumRangeLimit(limit, 1, 100)

	body := makeBody(ctx, utils.Map{"limit": limit})
	if v, ok := utils.StrEvaluation(sinceId, minId); ok {
		body["sinceId"] = v
	}
	if maxId != "" {
		body["untilId"] = maxId
	}
	_excludeTypes := lo.Map(excludeTypes,
		func(item models.NotificationType, _ int) models.MkNotificationType {
			return item.ToMkNotificationType()
		})
	_excludeTypes = append(_excludeTypes, models.MkNotificationTypeAchievementEarned)
	if lo.Contains(_excludeTypes, models.MkNotificationTypeMention) {
		_excludeTypes = append(_excludeTypes, models.MkNotificationTypeReply)
	}
	body["excludeTypes"] = _excludeTypes
	_includeTypes := lo.Map(types,
		func(item models.NotificationType, _ int) models.MkNotificationType {
			return item.ToMkNotificationType()
		})
	if lo.Contains(_includeTypes, models.MkNotificationTypeMention) {
		_includeTypes = append(_includeTypes, models.MkNotificationTypeReply)
	}
	if len(_includeTypes) > 0 {
		body["includeTypes"] = _includeTypes
	}

	var result []models.MkNotification
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/i/notifications"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	notifications := lo.Map(result, func(item models.MkNotification, _ int) models.Notification {
		n, err := item.ToNotification(ctx.ProxyServer())
		if err == nil {
			return n
		}
		return models.Notification{Type: models.NotificationTypeUnknown}
	})
	notifications = lo.Filter(notifications, func(item models.Notification, _ int) bool {
		return item.Type != models.NotificationTypeUnknown
	})
	return notifications, nil
}

func NotificationsClear(ctx Context) error {
	resp, err := client.R().
		SetBody(makeBody(ctx, nil)).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notifications/mark-all-as-read"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func NotificationsUnreadCount(ctx Context) (int, error) {
	var result struct {
		UnreadNotificationsCount int `json:"unreadNotificationsCount"`
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, nil)).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/i"))
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return 0, errors.WithStack(err)
	}
	return result.UnreadNotificationsCount, nil
}

func NotificationGet(ctx Context, id string) (models.Notification, error) {
	var mkNotification models.MkNotification
	body := makeBody(ctx, utils.Map{"notificationId": id})
	resp, err := client.R().
		SetBody(body).
		SetResult(&mkNotification).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notifications/show"))
	if err != nil {
		return models.Notification{}, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return models.Notification{}, errors.WithStack(err)
	}
	return mkNotification.ToNotification(ctx.ProxyServer())
}
