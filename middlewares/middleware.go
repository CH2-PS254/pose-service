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
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatError("Authentication credentials were missing or incorrect"))
		return
	}

	token, err := utils.ParseToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.FormatError("Failed to parse token"))
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("username", claims["username"])
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, utils.FormatError("The request is understood, but it has been refused or access is not allowed"))
		return
	}

	c.Next()
}
