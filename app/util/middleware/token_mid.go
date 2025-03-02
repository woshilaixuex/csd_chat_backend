package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/woshilaixuex/csd_chat_backend/app/util/security/xtoken"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-02-27 15:41
 */

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		token = token[7:]

		userID, err := xtoken.ParseJwtToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
