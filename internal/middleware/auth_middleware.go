package middleware

import (
	"net/http"
	"strings"

	"Blog_project_with_Go/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header topilmadi"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format noto'g'ri (Bearer <token> bo'lishi kerak)"})
			c.Abort()
			return
		}

		userID, err := utils.ValidateJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Kontekstga user_id ni saqlaymiz, keyingi handlerlarda ishlatish uchun
		c.Set("user_id", userID)
		c.Next()
	}
}