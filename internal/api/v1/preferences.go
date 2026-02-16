package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gizmo-ds/misstodon/models"
)

func PreferencesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, models.Preferences{
		PostingDefaultVisibility: "public",
		PostingDefaultSensitive:  false,
		ReadingExpandMedia:       "default",
		ReadingExpandSpoilers:    false,
	})
}
