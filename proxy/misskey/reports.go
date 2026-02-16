package misskey

import (
	"net/http"
	"strings"
	"time"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/pkg/errors"
)

func ReportCreate(ctx Context, accountID, comment string, statusIDs []string) (models.Report, error) {
	fullComment := comment
	if len(statusIDs) > 0 {
		fullComment += "\n\nRelated notes: " + strings.Join(statusIDs, ", ")
	}
	resp, err := client.R().
		SetBody(makeBody(ctx, utils.Map{
			"userId":  accountID,
			"comment": fullComment,
		})).
		Post(utils.JoinURL(ctx.ProxyServer(), "/api/users/report-abuse"))
	if err != nil {
		return models.Report{}, errors.WithStack(err)
	}
	if err = isucceed(resp, http.StatusNoContent); err != nil {
		return models.Report{}, errors.WithStack(err)
	}
	return models.Report{
		ID:          "0",
		ActionTaken: false,
		Category:    "other",
		Comment:     comment,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		StatusIDs:   statusIDs,
		RuleIDs:     []string{},
	}, nil
}
