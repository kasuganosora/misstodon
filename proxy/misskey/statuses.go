package misskey

import (
	"net/http"
	"time"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func StatusSingle(ctx Context, statusID string) (models.Status, error) {
	var status models.Status
	var note models.MkNote
	body := makeBody(ctx, utils.Map{"noteId": statusID})
	resp, err := client.R().
		SetBody(body).
		SetResult(&note).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/show"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, 200); err != nil {
		return status, errors.WithStack(err)
	}
	status = note.ToStatus(ctx.ProxyServer())
	if ctx.Token() != nil {
		state, err := getNoteState(ctx.ProxyServer(), *ctx.Token(), status.ID)
		if err != nil {
			return status, err
		}
		status.Bookmarked = state.IsFavorited
		status.Muted = state.IsMutedThread
	}
	return status, err
}

type noteState struct {
	IsFavorited   bool `json:"isFavorited"`
	IsMutedThread bool `json:"isMutedThread"`
}

func getNoteState(server, token, noteId string) (noteState, error) {
	var state noteState
	resp, err := client.R().
		SetBody(utils.Map{"i": token, "noteId": noteId}).
		SetResult(&state).
		Post(utils.JoinURL(server, "/api/notes/state"))
	if err != nil {
		return state, errors.WithStack(err)
	}
	if err = isucceed(resp, 200); err != nil {
		return state, errors.WithStack(err)
	}
	return state, nil
}

func StatusFavourite(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{
			"noteId":   id,
			"reaction": "⭐",
		})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/reactions/create"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "ALREADY_REACTED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Favourited = true
	status.FavouritesCount += 1
	return status, nil
}

func StatusUnFavourite(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/reactions/delete"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "NOT_REACTED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Favourited = false
	status.FavouritesCount -= 1
	return status, nil
}

func StatusBookmark(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/favorites/create"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "ALREADY_FAVORITED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Bookmarked = true
	return status, nil
}

func StatusUnBookmark(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/favorites/delete"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "NOT_FAVORITED"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Bookmarked = false
	return status, nil
}

// StatusBookmarks
// NOTE: 为了减少请求数量, 不支持 Bookmarked
func StatusBookmarks(ctx Context,
	limit int, sinceID, minID, maxID string) ([]models.Status, error) {
	var result []struct {
		ID        string        `json:"id"`
		CreatedAt string        `json:"createdAt"`
		Note      models.MkNote `json:"note"`
	}
	body := makeBody(ctx, utils.Map{"limit": limit})
	if v, ok := utils.StrEvaluation(sinceID, minID); ok {
		body["sinceId"] = v
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/i/favorites"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var status []models.Status
	for _, s := range result {
		status = append(status, s.Note.ToStatus(ctx.ProxyServer()))
	}
	return status, nil
}

// PostNewStatus 发送新的 Status
func PostNewStatus(ctx Context,
	status *string, pollOptions []string, pollExpiresIn int, pollMultiple bool,
	MediaIDs []string, InReplyToID string,
	Sensitive bool, SpoilerText string,
	Visibility models.StatusVisibility, Language string,
	ScheduledAt time.Time,
) (any, error) {
	body := makeBody(ctx, utils.Map{"localOnly": false})
	var noteMentions []string
	if status != nil && *status != "" {
		body["text"] = *status
		noteMentions = append(noteMentions, utils.GetMentions(*status)...)
	}
	if len(pollOptions) >= 2 {
		poll := utils.Map{
			"choices":  pollOptions,
			"multiple": pollMultiple,
		}
		if pollExpiresIn > 0 {
			poll["expiredAfter"] = pollExpiresIn * 1000 // Mastodon sends seconds, Misskey expects ms
		}
		body["poll"] = poll
	}
	if Sensitive {
		if SpoilerText != "" {
			body["cw"] = SpoilerText
		} else {
			body["cw"] = "Sensitive"
		}
	}
	switch Visibility {
	case models.StatusVisibilityPublic:
		body["visibility"] = "public"
	case models.StatusVisibilityUnlisted:
		body["visibility"] = "home"
	case models.StatusVisibilityPrivate:
		body["visibility"] = "followers"
	case models.StatusVisibilityDirect:
		body["visibility"] = "specified"
		var visibleUserIds []string
		for _, m := range noteMentions {
			a, err := AccountsLookup(ctx, m)
			if err != nil {
				return nil, err
			}
			visibleUserIds = append(visibleUserIds, a.ID)
		}
		if len(visibleUserIds) > 0 {
			body["visibleUserIds"] = visibleUserIds
		}
	}
	if MediaIDs != nil {
		body["mediaIds"] = MediaIDs
	}
	if InReplyToID != "" {
		body["replyId"] = InReplyToID
	}
	var result struct {
		CreatedNote models.MkNote `json:"createdNote"`
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/create"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	return result.CreatedNote.ToStatus(ctx.ProxyServer()), nil
}

func StatusDelete(ctx Context, id string) error {
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/delete"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func StatusReblog(ctx Context, id string) (models.Status, error) {
	original, err := StatusSingle(ctx, id)
	if err != nil {
		return models.Status{}, errors.WithStack(err)
	}
	var result struct {
		CreatedNote models.MkNote `json:"createdNote"`
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{
			"renoteId":   id,
			"localOnly":  false,
			"visibility": "public",
		})).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/create"))
	if err != nil {
		return models.Status{}, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return models.Status{}, errors.WithStack(err)
	}
	status := result.CreatedNote.ToStatus(ctx.ProxyServer())
	status.ReBlog = &original
	status.ReBlogged = true
	return status, nil
}

func StatusUnreblog(ctx Context, id string) (models.Status, error) {
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/unrenote"))
	if err != nil {
		return models.Status{}, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent); err != nil {
		return models.Status{}, errors.WithStack(err)
	}
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return models.Status{}, errors.WithStack(err)
	}
	status.ReBlogged = false
	return status, nil
}

func StatusRebloggedBy(ctx Context, id string, limit int, sinceID, maxID string) ([]models.Account, error) {
	var result []models.MkNote
	body := makeBody(ctx, utils.Map{"noteId": id, "limit": limit})
	if sinceID != "" {
		body["sinceId"] = sinceID
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/renotes"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var accounts []models.Account
	for _, note := range result {
		if note.User != nil {
			if a, err := note.User.ToAccount(ctx.ProxyServer()); err == nil {
				accounts = append(accounts, a)
			}
		}
	}
	return accounts, nil
}

func StatusFavouritedBy(ctx Context, id string, limit int, sinceID, maxID string) ([]models.Account, error) {
	type reactionResult struct {
		ID        string        `json:"id"`
		User      models.MkUser `json:"user"`
		Type      string        `json:"type"`
		CreatedAt string        `json:"createdAt"`
	}
	var result []reactionResult
	body := makeBody(ctx, utils.Map{"noteId": id, "limit": limit})
	if sinceID != "" {
		body["sinceId"] = sinceID
	}
	if maxID != "" {
		body["untilId"] = maxID
	}
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/reactions"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return nil, errors.WithStack(err)
	}
	var accounts []models.Account
	for _, r := range result {
		if a, err := r.User.ToAccount(ctx.ProxyServer()); err == nil {
			accounts = append(accounts, a)
		}
	}
	return accounts, nil
}

func StatusContext(ctx Context, id string) (map[string]any, error) {
	result := map[string]any{
		"ancestors":   []models.Status{},
		"descendants": []models.Status{},
	}

	// Get ancestors via /api/notes/conversation
	var ancestorNotes []models.MkNote
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id, "limit": 40})).
		SetResult(&ancestorNotes).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/conversation"))
	if err == nil && isucceed(resp, http.StatusOK) == nil {
		var ancestors []models.Status
		for _, n := range ancestorNotes {
			ancestors = append(ancestors, n.ToStatus(ctx.ProxyServer()))
		}
		// Reverse to chronological order
		for i, j := 0, len(ancestors)-1; i < j; i, j = i+1, j-1 {
			ancestors[i], ancestors[j] = ancestors[j], ancestors[i]
		}
		if len(ancestors) > 0 {
			result["ancestors"] = ancestors
		}
	}

	// Get descendants via /api/notes/children
	var childNotes []models.MkNote
	resp, err = client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id, "limit": 40})).
		SetResult(&childNotes).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/children"))
	if err == nil && isucceed(resp, http.StatusOK) == nil {
		var descendants []models.Status
		for _, n := range childNotes {
			descendants = append(descendants, n.ToStatus(ctx.ProxyServer()))
		}
		if len(descendants) > 0 {
			result["descendants"] = descendants
		}
	}

	return result, nil
}

func StatusMuteThread(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/thread-muting/create"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "ALREADY_MUTING"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Muted = true
	return status, nil
}

func StatusUnmuteThread(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/thread-muting/delete"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent, "NOT_MUTING"); err != nil {
		return status, errors.WithStack(err)
	}
	status.Muted = false
	return status, nil
}

func StatusPin(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/i/pin"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK, "ALREADY_PINNED"); err != nil {
		return status, errors.WithStack(err)
	}
	return status, nil
}

func StatusUnpin(ctx Context, id string) (models.Status, error) {
	status, err := StatusSingle(ctx, id)
	if err != nil {
		return status, errors.WithStack(err)
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{"noteId": id})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/i/unpin"))
	if err != nil {
		return status, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK, "NOT_PINNED"); err != nil {
		return status, errors.WithStack(err)
	}
	return status, nil
}

func StatusTranslate(ctx Context, id, targetLang string) (models.Translation, error) {
	var result struct {
		SourceLang string `json:"sourceLang"`
		Text       string `json:"text"`
	}
	body := makeBody(ctx, utils.Map{"noteId": id, "targetLang": targetLang})
	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/translate"))
	if err != nil {
		return models.Translation{}, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusOK); err != nil {
		return models.Translation{}, errors.WithStack(err)
	}
	return models.Translation{
		Content:                result.Text,
		DetectedSourceLanguage: result.SourceLang,
		Provider:               "Misskey",
	}, nil
}

func SearchStatusByHashtag(ctx Context,
	hashtag string,
	limit int, maxId, sinceId, minId string) ([]models.Status, error) {
	body := makeBody(ctx, utils.Map{"limit": limit})
	if v, ok := utils.StrEvaluation(sinceId, minId); ok {
		body["sinceId"] = v
	}
	if maxId != "" {
		body["untilId"] = maxId
	}
	body["tag"] = hashtag
	var result []models.MkNote
	_, err := client.R().
		SetBody(body).
		SetResult(&result).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/search-by-tag"))
	if err != nil {
		return nil, err
	}
	var list []models.Status
	for _, note := range result {
		list = append(list, note.ToStatus(ctx.ProxyServer()))
	}
	return list, nil
}
