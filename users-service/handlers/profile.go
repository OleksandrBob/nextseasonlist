package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProfileHandler(c *gin.Context) {
	userId, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user id": userId.(string)})
	return
}
