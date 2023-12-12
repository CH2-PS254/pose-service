package middlewares

import (
	"net/http"
	"pose-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Middleware(c *gin.Context) {
	tokenString, err := utils.ExtractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.ParseToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Malformed token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("username", claims["username"])
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
		return
	}

	c.Next()
}
