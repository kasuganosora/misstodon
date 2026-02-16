package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MarkersGetHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func MarkersPostHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
