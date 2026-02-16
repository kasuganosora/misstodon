package misskey

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func PollVote(ctx Context, noteID string, choices []int) (*models.Poll, error) {
	for _, choice := range choices {
		resp, err := client.R().
			SetBody(makeBody(ctx, utils.Map{"noteId": noteID, "choice": choice})).
			Post(utils.JoinURL(ctx.ProxyServer(), "/api/notes/polls/vote"))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err = isucceed(resp, http.StatusNoContent, "ALREADY_VOTED"); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	status, err := StatusSingle(ctx, noteID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return status.Poll, nil
}
