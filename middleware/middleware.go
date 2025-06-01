package middleware

import (
	"log"
	"net/http"
	"newsAPI/api"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CORSMiddleware настраивает CORS для API
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware проверяет JWT токен
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("AuthMiddleware: Processing request to %s", c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("AuthMiddleware: Authorization header is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Принимаем токен как есть, без проверки на Bearer
		tokenString := authHeader
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return api.GetJWTKey(), nil
		})

		if err != nil {
			log.Printf("AuthMiddleware: Token validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Printf("AuthMiddleware: Token is invalid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		log.Printf("AuthMiddleware: Successfully authenticated user ID: %v", claims.UserID)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
