package middleware

import (
	"net/http"
	"time"

	"github.com/coderunne/jwt-login/pkgs/models"
	"github.com/coderunne/jwt-login/pkgs/utility"
	"github.com/gin-gonic/gin"
)

func AuthenticateSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Verify and retrieve claims
		claims, err := utility.VerifyToken(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check and verify cookie expiry date
		now := float64(time.Now().Unix())
		if now > claims["expiry"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Retrieve user info from database using claims data
		user, err := models.GetUserByEmail(claims["issuer"].(string))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Store user data in context
		c.Set("user", user)

		// Continue Request
		c.Next()
	}
}
