package controllers

import (
	"net/http"
	"pose-service/db"
	"pose-service/models"
	"pose-service/utils"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var user models.User

	if err := db.GetDB().Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid credentials",
		})
		return
	}

	token, err := utils.GenerateToken(user.Username)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logged in successfully",
		"data": gin.H{
			"username": user.Username,
			"token":    token,
		},
	})
}
