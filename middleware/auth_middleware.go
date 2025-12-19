package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func AuthMiddlewareJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		userID, ok := claims["userID"].(string)
		if !ok || strings.TrimSpace(userID) == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
