package middlewares

import (
	"net/http"

	"github.com/OleksandrBob/nextseasonlist/users-service/utils"
	"github.com/gin-gonic/gin"
)

func BlacklistedTokensMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie(utils.RefreshTokenName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateRefreshToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is already invalid"})
			c.Abort()
			return
		}

		c.Set(utils.UserIdClaim, claims[utils.UserIdClaim])
	}
}
