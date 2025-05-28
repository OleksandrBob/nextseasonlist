package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"github.com/OleksandrBob/nextseasonlist/shared/token"
	"github.com/OleksandrBob/nextseasonlist/shared/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(accessTokecSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenStr := tokenParts[1]

		claims, err := token.ValidateAccessToken(tokenStr, accessTokecSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		userId, ok := claims[utils.UserIdClaim]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing userId in claims"})
			c.Abort()
			return
		}
		c.Set(utils.UserIdClaim, userId)

		userRoles, ok := claims[utils.RolesClaim]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing userId in claims"})
			c.Abort()
			return
		}
		c.Set(utils.RolesClaim, userRoles)

		c.Next()
	}
}

func AllowRoleMiddleware(allowedRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, claimsExist := c.Get(utils.RolesClaim)

		rolesIface, ok := userRoles.([]interface{})
		if !claimsExist || !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "User's role does not allow specified request"})
			c.Abort()
			return
		}

		roles := make([]string, len(rolesIface))
		for i, v := range rolesIface {
			roleStr, ok := v.(string)
			if !ok {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role format"})
				c.Abort()
				return
			}
			roles[i] = roleStr
		}

		if !slices.Contains(roles, allowedRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User's role does not allow specified request"})
			c.Abort()
			return
		}

		c.Next()
	}
}
