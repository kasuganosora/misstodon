package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ConversationsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, []any{})
}
